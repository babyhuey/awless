# AGENTS.md

Guide for AI agents working in this repository.

## Project Overview

`awless` is a CLI tool for managing AWS resources. It provides a template DSL for infrastructure creation/revert, local graph-based resource sync, smart SSH, and human-friendly output. This is a modernized fork of [wallix/awless](https://github.com/wallix/awless) maintained at [bootswithdefer/awless](https://github.com/bootswithdefer/awless).

- **Language:** Go 1.26 (`go 1.26.1` with `toolchain go1.26.5` in `go.mod`)
- **Module path:** `github.com/bootswithdefer/awless`
- **AWS SDK:** v2 (`github.com/aws/aws-sdk-go-v2`)
- **CLI framework:** `github.com/spf13/cobra`
- **Graph storage:** `github.com/bootswithdefer/triplestore` v1.0.0 (RDF triples)
- **Local DB:** bbolt (`go.etcd.io/bbolt`)
- **Version:** `v1.2.0` in `config/version.go`; release builds inject it from the git tag via GoReleaser. See [Releasing](#releasing).

## Build & Test

Use the Makefile rather than raw commands — it pins tool versions to match CI.
`make help` lists everything.

```sh
make build         # build the awless binary
make test          # go test ./...
make test-race     # go test -race ./...
make lint          # golangci-lint (pinned v2.12.2, installs on demand)
make vet           # go vet ./...
make vuln          # govulncheck ./...
make fmt           # gofmt -s + goimports -local github.com/bootswithdefer/awless
make fmt-check     # fail if anything non-generated is unformatted
make cover         # coverage profile + total (uses -coverpkg, see below)
make generate      # regenerate gen_*.go (see caveat below)
make fuzz          # run the native fuzzers briefly
make tools         # install pinned dev tools
make release-snapshot  # build release artifacts locally (no publish)

make check         # fast gate:  fmt-check vet lint test
make verify        # full gate:  fmt-check vet lint test-race vuln  (mirrors CI)
```

`make verify` is the gate to run before committing; it is what CI enforces.

**Check the exit status of `make verify` directly.** Piping it into `tail` or wrapping
it in `if ... then ... fi` gates on the wrong command's status and will happily let a
failing lint through.

`make cover` passes `-coverpkg=./...`. Without it, Go credits coverage only to the
package under test, and the large amount of `aws/spec` driven by the acceptance tests
is reported as uncovered — that alone understated the total by about seven points.

`make generate` requires `goimports` on PATH (`make tools` installs it): the
generators shell out to it to prune the unused imports their templates emit.
Output is deterministic, and CI enforces that committed generated files match
(the `codegen` job).

## Directory Structure

```
.
├── main.go                  # Entry point — the single os.Exit; see "Error handling"
├── commands/                # Cobra CLI commands (list, run, show, ssh, sync, etc.)
├── aws/
│   ├── services/            # Service implementations (gen_services.go is generated)
│   ├── spec/                # Command specs per resource (gen_runs.go, gen_cmds_defs.go, gen_inits.go generated)
│   ├── fetch/               # Data fetchers (gen_fetchers.go is generated)
│   ├── conv/                # AWS SDK type → internal model conversion
│   ├── config/              # AWS config validation
│   ├── doc/                 # CLI documentation: params, examples, enums
│   └── tailers/             # CloudFormation/ASG event tailers
├── cloud/                   # Cloud abstraction layer (interfaces, properties, RDF)
│   ├── properties/          # Generated property constants
│   └── rdf/                 # Generated RDF namespace constants
├── template/                # Template DSL engine
│   ├── internal/ast/        # PEG grammar + parser; gen_entities.go is generated
│   ├── env/                 # Template execution environment
│   └── params/              # Parameter validation
├── graph/                   # RDF-based resource graph
├── gen/aws/                 # Code generation
│   ├── generators/          # Generator programs (go run *.go)
│   ├── properties_definitions.go
│   ├── fetchers_definitions.go
│   └── mock_definitions.go
├── acceptance/aws/          # Acceptance test framework — 200 tests over all 194 commands
├── config/                  # App config, versioning, upgrade logic
├── console/                 # Terminal display, table formatting, column headers
├── database/                # bbolt-backed local storage
├── inspect/                 # Infrastructure analysis inspectors
├── ssh/                     # SSH client implementation
├── sync/                    # Cloud → local graph sync
├── web/                     # Web-based resource viewer
├── smoke_tests/             # Shell-based integration tests (require AWS credentials)
├── .github/workflows/ci.yml # GitHub Actions CI
└── .githooks/pre-commit     # gofmt + golangci-lint pre-commit hook
```

## Code Generation

Generated files follow the `gen_*.go` naming convention and are excluded from linting.

```sh
make generate   # or: cd gen/aws/generators && go run *.go
```

| Definition source | Output |
|----------------|--------|
| `gen/aws/properties_definitions.go` | `cloud/properties/gen_properties.go`, `cloud/rdf/gen_rdf.go` |
| `gen/aws/fetchers_definitions.go` | `aws/fetch/gen_fetchers.go` |
| `gen/aws/mock_definitions.go` | `aws/services/gen_mocks_test.go`, `acceptance/aws/gen_mocks.go` |
| `entity:` struct tags in `aws/spec/` | `template/internal/ast/gen_entities.go`, `acceptance/aws/gen_factory.go` |
| Definitions in `generators/*.go` | `aws/services/gen_services.go`, `aws/spec/gen_runs.go`, `aws/spec/gen_cmds_defs.go`, `aws/spec/gen_inits.go` |

**Do NOT edit `gen_*.go` files directly.** Edit the definitions in `gen/aws/` or the generator templates in `gen/aws/generators/`, then regenerate.

The generators validate that their output parses as Go before overwriting, so a broken
template fails without destroying the previous good file. They also apply an initialism
table (`capitalize` in `generators/main.go`) so the `dns` service generates `DNS` rather
than `Dns`. `api` is deliberately absent from that table, because `capitalize` also
renders SDK type names and apigatewayv2's type is `Api`.

## Adding a New AWS Service

Every step below is load-bearing; the ones
that fail *silently* are called out, because those are the ones that cost time.

```sh
go get github.com/aws/aws-sdk-go-v2/service/<svc>
```

**Definitions**

1. `cloud/cloud.go` — a resource type constant per resource. The value is what users type,
   so it must be a single lowercase word (`cachesubnetgroup`, not `cache-subnet-group`).
2. `gen/aws/properties_definitions.go` — any property not already there. The list is
   maintained alphabetically by `AwlessLabel`.
3. `gen/aws/fetchers_definitions.go` — one `fetchersDef` for the service. This drives
   **both** `gen_fetchers.go` and `gen_services.go`; there is nothing to add in
   `generators/services.go`. Get `NextPageMarker` from the SDK's output type rather than
   assuming `NextToken` — several services use `Marker`.
4. `aws/fetch/config.go` — add `<Svc> *<svc>.Client` to `AWSAPI`. `assignAPIs` matches on
   type by reflection, so no explicit assignment is needed anywhere.
5. `aws/fetch/manual_fetchers.go` — add `func addManual<Svc>FetchFuncs(conf *Config, funcs
   map[string]fetch.Func) {}`, empty if the service needs no manual fetchers. The generated
   code calls it unconditionally, so its absence is a compile error.
6. `make generate`.

**Wiring**

7. `aws/conv/convert.go` — a `case` per AWS type, plus the `<svc>types` import alias.
   `aws/conv/model.go` — the property mapping. Use `extractValueFn` for a plain `[]string`;
   `extractStringSliceValues("Field")` is for slices *of structs*.
8. `aws/services/init.go` — declare the var, call `New<Svc>`, and add the
   `cloud.ServiceRegistry` entry. A service that is generated but never registered compiles
   fine and is **invisible at runtime**. `TestEveryGeneratedServiceIsRegistered` catches
   this now.
9. `aws/spec/<svc>.go` — the command specs. See
   [Command Specs and the Reflective Setters](#command-specs-and-the-reflective-setters);
   **verify every `awsName` against the SDK**, because one that does not resolve is dropped
   without a word.
10. `console/defaults.go` — **two** places: the short property list near the top, and the
    full `ColumnDefinition` block lower down. Nothing is needed in `console/headers.go`,
    which is generic machinery rather than a per-resource registry.
11. `template/revert.go` — if the resource's `delete` takes `name` rather than `id`, add the
    entity to the `name=` case in the `create` branch. The default emits `id=<result>`, and
    a revert that does not match its own delete command fails with `unexpected param id`.
    Eight resources shipped with this broken.
12. `aws/doc/paramsdoc.go` and `aws/doc/clidoc.go` — a doc line per param and at least one
    worked example per command. `TestDocForEachParam` and `TestDocForEachCommand` fail
    without them, and `TestExamplesSatisfyParamsSpec` validates each example against the
    spec.

**Tests and docs**

13. `acceptance/aws/<svc>_test.go` — assert the **input mapping**, not just that the call
    happened. Include an empty-response case: an AWS reply with a nil body is the shape that
    made six commands panic during the SDK v2 migration.
14. `aws/services/services_registry_test.go` — add the service to the explicit list.
15. `README.md` service table and `ls` examples, and a `CHANGELOG.md` entry.

### Conventions that are not obvious

- **Integer params are `*int64` with `awsType:"awsint64"`** even when the AWS field is
  `*int32`. The setter converts widths through `reflect.Convert`; there is no `awsint32`.
- **Shared params go above a `params.OnlyOneOf`, not inside each branch.** The param
  collector walks every branch, so anything repeated is listed twice in `-h`. The trade-off
  is that a param only meaningful to one form is accepted and then rejected by AWS.
- **`ExtractResult` should fall back to a param** when the response body may be empty,
  rather than dereferencing it.
- **`RunExpectingError` returns the error**, and `errcheck` requires it to be checked.
- Verify the `awsName` tags mechanically before committing. A `go doc` of each `awsInput`
  type compared against the tags in the spec catches the whole class in one pass.

The template parser's entity list is generated from the `entity:` struct tags, so a new
entity needs no manual registration. That was not always true: a command could
previously be fully registered and still fail every template with `unknown entity`.

## Command Specs and the Reflective Setters

A command is a struct in `aws/spec/` whose tags drive everything:

```go
type CreateVpc struct {
    _    string  `action:"create" entity:"vpc" awsAPI:"ec2" awsCall:"CreateVpc" awsInput:"ec2.CreateVpcInput" awsOutput:"ec2.CreateVpcOutput" awsDryRun:""`
    CIDR *string `awsName:"CidrBlock" awsType:"awsstr" templateName:"cidr"`
}
```

`templateName` is what users type. `awsName` is the field path on the AWS input, and
`awsType` selects a conversion in `aws/spec/setters.go`. **All of this is applied by
reflection, so the compiler checks none of it.** Nine bugs of exactly this shape were
found and fixed; the recurring causes are worth knowing:

- **SDK v1 vs v2 shapes.** v1 modeled lists as `[]*string` and `[]*Struct`; v2 uses
  `[]string` and `[]Struct`. The setters convert both ways now, but a new `awsType` must
  handle the value form.
- **A wrong `awsName` is silent.** `setValueAtPath` ignores a field it cannot find,
  deliberately, since these tags are usually generated. A case mismatch
  (`Healthcheck` vs `HealthCheck`) therefore made a command send an empty request and
  report success.
- **A wrong `awsType` is loud.** `setFieldWithType` returns an error naming the tag and
  field for an unrecognized type. That guard caught an `awsin64` typo.
- **Field paths may be indexed**, as in `DistributionConfig.Origins.Items[0].DomainName`.

`awsCall` commands with a hand-written `ManualRun` build their own request; those use
`renv.RequestContext()` for the AWS call.

## Error Handling and Exit Codes

`main.go` holds the **single** `os.Exit`. Commands return errors through cobra, so
deferred cleanup always runs, and `RootCmd` sets `SilenceErrors` so reporting happens in
exactly one place. Two sentinels in `commands/errors.go` carry outcomes that are not
plain failures:

- `ErrExitZero` — the command has already told the user what they need; exit 0 silently.
- `ErrReported` — the reason was already printed, usually with a suggested command; exit
  non-zero without printing again.

Two `os.Exit` calls remain outside `main`, both readline interrupt paths whose enclosing
signature has no error to return (a template env `MissingHolesFunc` and a
`stdinParamProviderFn`). Both close readline explicitly first, because `os.Exit` skips
the deferred close and would leave the terminal in raw mode. Both record that reasoning
in place.

## Context

`main.go` installs a `signal.NotifyContext` root context, reachable as
`commands.RootContext()` and threaded to commands as `env.Running.RequestContext()`.
`Context()` on that interface is the **template variable map**, not a `context.Context` —
they are easy to confuse.

Every outbound call must carry the context: `noctx` and `contextcheck` are enabled and
have each caught commands that could not be interrupted.

## Acceptance Tests

`acceptance/aws/` drives real command execution with no network access. SDK v2 exposes
concrete `*service.Client` structs, so mocking works by injecting a smithy middleware on
the Initialize step — the outermost step, which is the only one still carrying typed
input parameters and whose short-circuit result becomes the operation output.

```go
mock := NewMock().On("CreateVpc", &ec2.CreateVpcOutput{Vpc: &ec2types.Vpc{VpcId: awssdk.String("vpc-1234")}})

Template("create vpc cidr=10.0.0.0/16").
    Mock(mock).
    ExpectCalls("CreateVpc").
    ExpectCommandResult("vpc-1234").
    ExpectRevert("delete vpc id=vpc-1234").
    Run(t)

in := mock.InputFor("CreateVpc").(*ec2.CreateVpcInput)   // assert the input mapping
```

Also available: `RunExpectingError` for asserting pre-flight validation, `DryRun()` to
exercise the generated `dryRun` path, and `OnAPIError` for AWS error codes. An operation
with no registered output returns an explanatory error rather than a zero value, since a
zero value surfaces as a nil dereference deep inside result extraction.

Assert the **input mapping**, not just that a call happened — that is what catches the
reflective-tag bugs. If a test needs an unexpected extra mock, read why: several
commands legitimately make more calls than they appear to (`delete role` also tears down
its instance profile, `create securitygroup`'s revert waits for the group to become
unused).

Keep timeouts small in `check *` tests. Those commands poll, so a fixture whose state
never matches hangs the suite instead of failing it.

## Conventions

- **Imports:** stdlib, then third-party, then `github.com/bootswithdefer/awless` (enforced by goimports with `-local`)
- **Error handling:** wrap with `%w`; use `decorateAWSError()` for AWS errors
- **Pointer helpers:** `String()`, `StringValue()`, `Int64()`, `Bool()` from `aws/spec/spec.go`
- **Resource types:** constants in `cloud/cloud.go` (e.g. `cloud.Instance`, `cloud.Vpc`)
- **Spelling:** American. `misspell` runs with `locale: US`, so `behaviour` and
  `initialise` fail the build.
- **Generated file prefix:** `gen_` — never hand-edit
- **Test file suffix:** `_test.go`, or `_extra_test.go` for external test packages

## Linting

`.golangci.yml` enables 51 linters. The set was expanded deliberately, and the config
records why anything is absent — `dupword` and `unparam` were run once, found only false
positives or a legitimate builder pattern, and are excluded with that reasoning in place.

`SA1019` (deprecated usage) is suppressed; `gen_*.go` and the PEG parser are excluded;
`fieldalignment` and `shadow` are off for govet, and so is `inline`.

`gosec` is deliberately **not** enabled — see [Deliberate Omissions](#deliberate-omissions). It was reviewed once
rather than adopted: over half its findings duplicate `errcheck` or describe inherent
properties of a CLI that reads user-specified files and URLs. Worth re-running
periodically rather than gating on.

## CI

`.github/workflows/ci.yml` runs `test`, `race`, `codegen`, `fuzz`, `vuln`, `lint` and
`build`. Actions are pinned to commit SHAs (update with `make pinact-update`), and the Go
version comes from `go.mod`.

CI runs on every push to `master` and on pull requests, and is green: all twelve jobs
pass, including the codegen drift check, the race detector, `govulncheck`, a fuzz smoke
run and cross-compilation for six platforms.

Note that CI is a stronger check than local `make verify`, because it runs on a clean
runner with no `~/.awless`, no AWS config and no cached tools. Two tests needed
environment pinned for that reason — the service registry test sets static credentials to
avoid a 15-second EC2 metadata timeout, and `create keypair` needs `__AWLESS_KEYS_DIR`.

## Releasing

A release is cut by pushing a `v*` tag. `.github/workflows/release.yml` then runs
GoReleaser, which builds six platform archives, writes `checksums.txt`, publishes the
GitHub release, and commits a Homebrew cask to `bootswithdefer/homebrew-tap`.

**Once a version is on the module proxy it is immutable.** `proxy.golang.org` caches what
a tag pointed at, so a mistake cannot be corrected by moving the tag — it needs a new
patch version. Do the pre-flight.

```sh
# 1. Bump the fallback constant and add the changelog entry.
#    config.Version is only used by a plain `go build`; GoReleaser injects the real
#    version from the tag, so this is hygiene rather than the source of truth.
$EDITOR config/version.go CHANGELOG.md

# 2. Local gate.
make verify

# 3. Commit, push, and wait for CI to go green on that commit.
#    CI is the stronger check — clean runner, no ~/.awless, no AWS config, no cached
#    tools. Do not tag a commit whose CI has not finished.
gh api "repos/bootswithdefer/awless/actions/runs?per_page=1" \
  --jq '.workflow_runs[0] | "\(.head_sha[0:8]) \(.status)/\(.conclusion)"'

# 4. Signed, annotated tag. Commits here are signed, so tags should be too.
git tag -s v1.1.1 -F -   # summary of the release in the message body
git tag -v v1.1.1        # expect "Good signature"
git push origin v1.1.1
```

Verify afterwards that the release has its six archives plus checksums, and that the cask
landed in the tap:

```sh
gh api repos/bootswithdefer/awless/releases --jq '.[0] | "\(.tag_name) assets=\(.assets|length)"'
gh api repos/bootswithdefer/homebrew-tap/contents/Casks --jq '.[].name'
```

### Version numbering

Semver here is relative to **this module's** history, not to upstream. The module path
changed to `github.com/bootswithdefer/awless`, so the fork's first release was `v1.1.0`
despite carrying breaking changes against `wallix/awless` — there was no prior release of
this module to break.

**Do not tag `v2.0.0` without renaming the module.** Go requires a major version of two or
above to appear in the module path, so `v2` would mean moving to
`github.com/bootswithdefer/awless/v2`, rewriting every import and changing the documented
install command. A module without the suffix is rejected outright:
`version "v2.0.0" invalid: should be v0 or v1, not v2`. That suffix exists so importers can
depend on two majors at once, which is meaningless for a CLI nothing imports.

### Homebrew

The cask is committed over SSH with a **deploy key** scoped to the tap, held in the
`HOMEBREW_TAP_SSH_KEY` secret. The workflow writes it to `~/.ssh/homebrew_tap` because
GoReleaser wants a path rather than key contents.

A personal access token would be the obvious alternative and is the wrong choice: the
workflow's `GITHUB_TOKEN` cannot reach another repository, and the narrowest classic PAT
that can needs `repo` scope, which grants write access to every repository the owner has.

It is a cask rather than a formula because it ships a pre-built binary, which is what
Homebrew expects casks to carry. GoReleaser's `brews` section is deprecated for the same
reason — use `homebrew_casks`. Get its field names from `goreleaser schema` rather than
from memory; several are deprecated within it too.

## Deliberate Omissions

Things that look like oversights and are not. Recorded here so they are not "fixed" by
someone who assumes nobody looked.

**`gosec` is not enabled as a gate.** It was run once over the whole tree as a review, and
produced 127 findings: 57 duplicate `errcheck`, which is already on; 14 are `G304` file
paths from variables, inherent to a CLI that takes `--template-file`; 11 are `0644` writes
already reviewed; 10 are false positives in doc strings and fixtures; 4 are `docker` and
`ssh` invocation. Its `G115` integer-overflow rule did find two real crashes in the network
monitor — a divide-by-zero when every request completes in the same instant, and an
unsigned wrap on a narrow terminal — both fixed. Worth re-running periodically; not worth
blocking a build on.

**`math/rand` is used unseeded** in `aws/spec/spec.go` and `graph/rdf.go`. Both sites
generate dry-run IDs and graph identifiers where unpredictability does not matter, and the
global source has been auto-seeded since Go 1.20. Left alone deliberately rather than
switched to `crypto/rand`.

**Classic ELB (v1) is still supported** alongside ELBv2, so the tree depends on both SDK
modules. EC2-Classic was retired in 2022, but Classic Load Balancers still exist in older
VPC accounts. Kept so the maintenance cost is a choice; dropping it would remove one SDK
module and `aws/spec/classicloadbalancer.go`.

**golangci-lint is installed with a pinned `go install`** rather than
`golangci-lint-action`. The action is faster, but pinning the exact version (v2.12.2) in the
Makefile is what keeps CI and local runs from disagreeing — an unpinned `@latest` is what
broke the lint job once already.

**Some AWS resources are deliberately not covered.** Every service that was on the
candidate list has been added; these specific resources within them were left out, and the
reason matters because each looks like an omission:

| Resource | Why |
|---|---|
| Step Functions and CodePipeline executions, CodeDeploy deployment history | Run history rather than infrastructure, and listing costs one call per parent |
| Global Accelerator endpoint groups | A third level down, behind a listener, with a document-shaped configuration |
| FSx volumes | Only exist for two of the four file system types, and list per file system |
| MSK and Amazon MQ configurations | Versioned blobs of engine properties rather than infrastructure |
| Cloud Map instances | Registered by whatever runs the workload, usually ECS, rather than by hand |
| AWS Backup recovery points | The backups themselves; a delete command here exists only to destroy data |
| VPC endpoint services | The provider side of PrivateLink, a different job from consuming an endpoint |
| WAF web ACL and rule group associations | Attaching an ACL to a load balancer is a separate API from managing the ACL |

The general rule these follow: awless syncs a graph of infrastructure, so a resource whose
population scales with *activity* rather than with what you have built does not belong in it.
Adding any of them is a normal service addition — the procedure above covers it.

**Five functions are deliberately untested.** `InteractiveTerminal`, `propagateSignals`,
`NewClientWithProxy`, `Connect` and `workaroundExeCVEThroughScript` need a live SSH
handshake or replace the process through `syscall.Exec`. Covering them means an
integration test with an in-process SSH server, not a seam. Everything else in the
terminal, stdin and network paths now has one — see the seam table below.

## Test Seams

Three seams exist purely for testability, each defaulting to the real implementation so
production behavior is unchanged. All are unexported or package-level, and tests restore
them with `t.Cleanup`.

| Seam | Falls back to | Why |
|---|---|---|
| `ssh.Client.dialFunc` | `gossh.Dial` via `dialer()` | exercise username iteration and partial failure without a network |
| `awsconfig.promptStdin` | `os.Stdin` | drive the region and instance-type prompts without a terminal |
| `console.termGetSize` | `term.GetSize(os.Stdout)` | stdout is a pipe under `go test`, so the real call fails |

`dialer()` is a fallback rather than a constructor default because `Client` is built
directly in several places; a zero-value `Client` must still dial for real.

Two test details that are easy to lose:

- The prompt-selector tests run under a 5s watchdog goroutine. A regression there is an
  infinite loop, not a wrong value, so without the watchdog the suite hangs instead of
  failing. Both selectors previously spun forever on a closed stdin.
- Tests that touch region resolution must set `AWS_EC2_METADATA_DISABLED`, or they wait out
  the EC2 instance metadata timeout — 5s instead of 0.01s per test.

## Things to Watch Out For

- **`triplestore`'s `(*triple).key()` format is a wire contract.** Changing it
  invalidates every binary-encoded graph already cached in users' `~/.awless`.
- **Persisted template log lines are re-parsed** by `UnmarshalJSON`, so redaction
  placeholders must stay grammar-parseable (`*****` is in the `UnquotedParam` class).
- **`template/env` exports were renamed** from `ALL_CAPS` (`env.FILLERS` is now
  `env.Fillers`).
- **`wallix/awless-templates` references in `commands/run.go`** point at a *different*
  upstream repo and must stay.
- **Template PEG regeneration** requires the `peg` tool
  (`go install github.com/pointlander/peg@latest`). Edit the `.peg` file, never the
  `.peg.go`.
- **A bare UUID is not a valid param value.** The grammar accepts `name=build-and-deploy`
  and `execution=abcd-1234`, but rejects `execution=12345678-1234-1234-1234-123456789012` —
  an all-numeric hyphenated value appears to parse as an expression. Quote it, and say so in
  the param documentation for any command that takes one.
- **Inline JSON must be single-quoted.** There is no escape for a double quote inside a
  double-quoted value, so `pattern='{"source":["aws.ec2"]}'` is the only form that works.
  For anything longer than a line, take a file instead — see `create statemachine`.
- **The scheduler is gone.** `--run-in` and `--revert-in` were removed in v1.1.0.
