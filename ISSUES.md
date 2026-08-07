# ISSUES.md

Known bugs, technical debt, and improvement opportunities in this repository.

---

## Bugs

### B1: Compiled binaries committed to git — **FIXED**

**Severity:** Medium  
**Files:** `./generators` (6.1MB), `gen/aws/generators/generators` (5.3MB)

Two ELF binaries are committed to the repository. They bloat clone size and are platform-specific (linux/amd64 only). The `.gitignore` excludes only the `awless` binary but not the generators.

**Fix:** Add `generators` and `gen/aws/generators/generators` to `.gitignore`, remove them from tracking with `git rm --cached`.

---

### B2: Upgrade check points to upstream wallix/awless

**Severity:** Low  
**File:** `config/upgrade.go:56`

The `notifyIfUpgrade` function fetches version info and prints upgrade instructions referencing `github.com/wallix/awless` — the original upstream repo, not this fork. Users following those instructions would install the wrong binary.

**Fix:** Update URLs and `go install` paths to reference `github.com/babyhuey/awless`.

---

### B3: `math/rand` used without seeding (pre-Go 1.20 pattern)

**Severity:** Low  
**Files:** `aws/spec/spec.go:6`, `graph/rdf.go:5`

Uses `math/rand` for dry-run ID generation and graph operations. Since Go 1.20, the global rand source is auto-seeded, so this is not technically broken on the current Go 1.26, but the pattern is outdated and could produce predictable values on older compilers.

**Fix:** Replace with `crypto/rand` where uniqueness matters, or leave as-is given Go 1.26 auto-seeds.

---

### B4: CI actions pinned to tags instead of commit SHAs — **FIXED**

**Severity:** Low  
**File:** `.github/workflows/ci.yml`

Actions are referenced by mutable tags (`actions/checkout@v4`, `actions/setup-go@v5`, `ncruces/go-coverage-report@v0`). A compromised tag could inject malicious code into the CI pipeline.

**Fix:** Pin to full commit SHAs (e.g., `actions/checkout@<sha>`).

---

### B5: golangci-lint config incompatible with v2 — **FIXED**

**Severity:** High  
**Files:** `.golangci.yml`, `.github/workflows/ci.yml`

The CI workflow installs golangci-lint with `go install ... @latest`, which now resolves to v2.12+. The `.golangci.yml` uses v1 config format (missing `version: "2"` field, uses `disable-all`, `linters-settings` at top level, `issues.exclude-dirs`, `issues.exclude-files`). Running it produces:

```
Error: can't load config: unsupported version of the configuration
```

Key v1→v2 breaking changes affecting this config:
- `linters.disable-all: true` → `linters.default: none`
- `linters-settings` → split into `linters.settings` and `formatters.settings`
- `gofmt` and `goimports` moved from linters to `formatters.enable`
- `gosimple` merged into `staticcheck`
- `issues.exclude-dirs` → `linters.exclusions.paths`
- `issues.exclude-files` → `linters.exclusions.paths`
- `goimports.local-prefixes` → `formatters.settings.goimports.local-prefixes` (with `sections` syntax)
- Must add `version: "2"` at top level

**Fix:** Migrate `.golangci.yml` to v2 format. The equivalent v2 config:

```yaml
version: "2"

linters:
  default: none
  enable:
    - govet
    - ineffassign
    - staticcheck
    - unused
    - misspell
    - unconvert
  settings:
    govet:
      enable-all: true
      disable:
        - fieldalignment
        - shadow
    misspell:
      locale: US
    staticcheck:
      checks:
        - all
        - -SA1019
  exclusions:
    generated: strict
    paths:
      - vendor
      - "gen_.*\\.go$"
      - "awless-template-syntax\\.peg\\.go$"

formatters:
  enable:
    - gofmt
    - goimports
  settings:
    goimports:
      local-prefixes:
        - github.com/wallix/awless

issues:
  max-issues-per-linter: 50
  max-same-issues: 5
```

**Also fix:** Pin golangci-lint version in CI (either use the official action `golangci/golangci-lint-action@v6` with a pinned version, or `go install ... @v2.12.2`). The `@latest` pattern causes breakage on every major release.

---

### B6: Secrets persisted in plaintext to the local template log — **FIXED**

**Severity:** High  
**Files:** `template/marshal.go:69`, `template/marshal.go:61`, `database/templates.go:41`

`TemplateExecution.MarshalJSON` serializes each executed command via `newCmd.Line = cmd.String()`, and stores `out.Fillers` verbatim. Every template execution is then persisted to BoltDB (`database/templates.go:41`).

Commands that take secrets as parameters therefore write those secrets to disk in plaintext, and `awless log` will print them back:

- `create loginprofile username=X password=SECRET` (`aws/spec/loginprofile.go:33`)
- `update loginprofile username=X password=SECRET` (`aws/spec/loginprofile.go:53`)
- `create database ... password=SECRET` (`aws/spec/database.go`)
- Any template `Fillers` holding sensitive values

The DB file itself is mode `0600` (`database/db.go:69`), so this is not world-readable, but secrets are retained indefinitely with no redaction or expiry, and are exposed by any backup, sync, or `awless log` output shared for support.

**Fix:** Mark sensitive parameters in the command struct tags (e.g., `sensitive:"true"`) and redact them in `cmd.String()` output and in `Fillers` before marshalling. Replace values with `<redacted>`.

---

### B7: `awless list --sort` panics on mixed-type columns — **FIXED**

**Severity:** Medium  
**File:** `console/displayer.go:975`, `console/displayer.go:1001`

The sort comparator panics rather than returning an error:

```go
if reflect.TypeOf(a) != reflect.TypeOf(b) {
    panic(fmt.Sprintf("can not compare values of type %T and %T", a, b))
}
```

and in the `default` branch of the type switch for any unhandled type. Since resource properties are `interface{}` sourced from heterogeneous AWS responses, a column whose value is a string on one resource and a number (or nil-typed differently) on another will crash the CLI mid-render.

Note `git log` shows a prior "sort panic" fix (`71665081`), so this class of crash has already been hit once in practice.

**Fix:** Return a stable ordering instead of panicking — fall back to `fmt.Sprint(a) <= fmt.Sprint(b)` for mismatched or unknown types. Panicking in display code is never the right failure mode for a CLI.

---

### B8: Goroutine leak on first error in fan-out helpers — **FIXED**

**Severity:** Medium  
**Files:** `aws/fetch/s3_helpers.go`, `aws/fetch/ecs_helpers.go`, `aws/fetch/manual_fetchers.go`, `aws/conv/convert.go`, `inspect/inspectors/pricer.go`

The pattern spawned one goroutine per item writing to an **unbuffered** error
channel while the consumer returned on the first error, leaving every other
failing goroutine blocked on send forever, so `wg.Wait()` never completed.

All 10 per-item fan-out sites are now `errgroup` with `SetLimit`, which fixes the
leak and bounds concurrency together.

Three sites had defects beyond the leak:

- `ecs_helpers.go getAllTasks` called `tasksWG.Add` from inside a worker while
  another goroutine could already be in `tasksWG.Wait()`, which is undefined
  behavior, and sent a result after sending an error.
- `manual_fetchers.go` IAM users registered `defer wg.Done()` **after** an early
  return, so a failing `InitResource` left the counter high and `wg.Wait()`
  deadlocked. It also wrote a shared `hasError` flag from several goroutines
  without synchronization.
- `aws/conv/convert.go` continued to send a result after sending an error.

The remaining argument-less `go func()` calls in `manual_fetchers.go` are
fixed-count producer goroutines (two or three per fetch), not per-item fan-outs,
so they do not scale with account size. They are a different and much lower-risk
pattern.

---

### B9: Web UI binds to all interfaces with no timeouts or auth — **FIXED**

**Severity:** Medium  
**Files:** `web/web.go:46`, `commands/web.go:35`

```go
return http.ListenAndServe(s.port, s.routes())
```

Three problems:

1. **Default port `:8080` binds `0.0.0.0`**, not `127.0.0.1` — the local infrastructure graph is exposed to the whole network by default.
2. **No authentication** on any route (`/resources`, `/rdf`, `/graph`, `/`). It serves the synced AWS inventory: instance IDs, private IPs, security-group rules, VPC topology.
3. **No server timeouts.** Bare `http.ListenAndServe` has no `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, or `IdleTimeout`, making it Slowloris-exposed.

**Fix:** Default the bind address to `127.0.0.1:8080`, and construct an explicit server:

```go
srv := &http.Server{
    Addr:              s.port,
    Handler:           s.routes(),
    ReadHeaderTimeout: 5 * time.Second,
    ReadTimeout:       30 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       120 * time.Second,
}
return srv.ListenAndServe()
```

If non-loopback binding is ever wanted, require an explicit opt-in flag and warn that the endpoint is unauthenticated.

---

### B10: `panic()` in library code paths

**Severity:** Low  
**Files:** `sync/sync.go:61`, `sync/repo/repo.go:106`, `graph/graph.go:387`

- `NewSyncer` panics if the repo cannot be created (e.g., permissions on `~/.awless`) instead of returning an error, even though the surrounding API is error-returning.
- `repo.go:106` panics while listing revisions.
- `Graph.MustMarshal` panics by design, but a non-panicking `MarshalTo` already exists next to it — callers should prefer it.

**Fix:** Change `NewSyncer` to return `(Syncer, error)`. Audit `MustMarshal` call sites and migrate to `MarshalTo`.

---

### B11: Generated source written world-writable (`0666`) — **FIXED**

**Severity:** Low  
**File:** `gen/aws/generators/main.go:75`

```go
if err := os.WriteFile(path, buff.Bytes(), 0666); err != nil {
```

Generated Go source files are written mode `0666` (world-writable, subject to umask). Every other write in the codebase correctly uses `0600`/`0700`. On a shared host with a permissive umask, another user could modify generated source that then gets compiled into the binary.

**Fix:** Change to `0644`.

---

## Technical Debt

### D1: Acceptance tests completely broken

**Severity:** High  
**File:** `acceptance/aws/gen_factory.go`

All 30 factory functions contain `// TODO: SDK v2 mocking needs rework`. The acceptance test framework cannot inject mocked API clients since the SDK v2 migration changed from interface-based to concrete `*service.Client` types.

**Fix:** Implement interface wrappers for each SDK client (or use `smithy-go` middleware mocking) and regenerate the factory. This is a significant refactor.

---

### D2: `strings.Title` deprecated since Go 1.18 — **FIXED**

**Severity:** Medium  
**Files:**
- `gen/aws/generators/mocks.go:38`
- `gen/aws/generators/services.go:38`
- `gen/aws/generators/acceptance_mocks.go:42,68`
- `gen/aws/generators/fetchers.go:28`
- `gen/aws/fetchers_definitions.go:46`
- `gen/aws/mock_definitions.go:181`
- `console/displayer.go:82,252`
- `graph/resource.go:374`
- `aws/spec/policy.go:347`
- `commands/run.go:327`

`strings.Title` was deprecated in Go 1.18 because it doesn't handle Unicode word boundaries correctly. The lint config suppresses `SA1019` to avoid noise, but these should be migrated.

**Fix:** Replace with `golang.org/x/text/cases` or a simple `strings.ToUpper(s[:1]) + s[1:]` where input is ASCII-only (which it is in all generator cases).

---

### D3: Stale `.travis.yml` — **FIXED**

**Severity:** Low  
**File:** `.travis.yml`

References Go 1.9–1.11 and `goveralls`. Completely superseded by `.github/workflows/ci.yml`. Misleading for contributors.

**Fix:** Delete the file.

---

### D4: 101 uses of `context.Background()` in non-test code

**Severity:** Medium  
**Files:** `aws/spec/gen_runs.go` (generated), `aws/fetch/manual_fetchers.go`, `aws/tailers/`, `sync/sync.go`

AWS SDK v2 calls use `context.Background()` everywhere instead of accepting a context parameter. This prevents proper cancellation, timeout propagation, and graceful shutdown.

**Fix:** Thread `context.Context` through the command and fetcher interfaces. This is a large cross-cutting change since `gen_runs.go` is generated from templates in `gen/aws/generators/commands.go`.

---

### D5: New services are list-only (no CRUD commands)

**Severity:** Medium  
**Files:** New services (EKS, DynamoDB, Secrets Manager, API Gateway, SSM, EFS, CloudTrail, CloudWatch Logs) in `aws/services/gen_services.go`

The 8 newly-added services only support fetching/listing. They have no create/update/delete command specs in `aws/spec/`. Users can `awless list eksclusters` but cannot `awless create ekscluster`.

**Fix:** Add command specs for common CRUD operations per service. Each requires a new file in `aws/spec/` following the existing pattern (e.g., `instance.go`).

---

### D6: Deprecated / abandoned dependencies — **FIXED**

**Severity:** Medium  
**File:** `go.mod`

Verified against upstream repositories. No module declares a formal `Deprecated:` directive in its `go.mod`, so `go list -m -u` will **not** warn about any of these — they must be tracked manually.

| Module | Current | Status | Replacement |
|--------|---------|--------|-------------|
| `github.com/boltdb/bolt` | v1.3.1 | **Archived Mar 9, 2019** | `go.etcd.io/bbolt` v1.5.0 |
| `gopkg.in/src-d/go-git.v4` | v4.13.1 | **Archived Sep 11, 2020** (src-d shut down) | `github.com/go-git/go-git/v5` v5.19.2 |
| `gopkg.in/yaml.v2` | v2.4.0 | **Archived Apr 1, 2025** | `go.yaml.in/yaml/v3` v3.0.5 |
| `github.com/oklog/ulid` | v1.3.1 | v1 superseded | `github.com/oklog/ulid/v2` v2.1.2 |
| `github.com/wallix/triplestore` | 2018-02-13 pseudo-version | No tagged release, untouched since Feb 2018; same abandoned upstream org as awless itself | No successor — fork or vendor |
| `github.com/wallix/awless-scheduler` | v0.0.6 | Same abandoned upstream org | No successor — fork or vendor |
| `github.com/mitchellh/ioprogress` | 2018-02-01 pseudo-version | No tagged release, untouched since Feb 2018 | Small enough to inline (~100 LOC) |

**Note on yaml:** both `yaml.v2` **and** `yaml.v3` are unmaintained — `go-yaml/yaml` was archived entirely in Apr 2025. The maintained successor is the YAML-org fork published as `go.yaml.in/yaml/v3`. Migrating v2 → v3 alone is therefore not sufficient. `go.yaml.in/yaml/v3` is *already* in the dependency tree transitively, so this consolidates rather than adds a dependency.

**Usage scope (effort per module):**

| Module | Usage sites |
|--------|-------------|
| `yaml.v2` | 1 file — `aws/spec/stack.go` |
| `oklog/ulid` | 2 files — `template/marshal.go`, `template/template.go` |
| `ioprogress` | 1 file — `aws/spec/s3object.go` |
| `boltdb/bolt` | 2 files — `database/db.go`, `database/templates.go` |
| `src-d/go-git.v4` | 1 file — `sync/repo/repo.go` |
| `triplestore` | Core — `graph/`, `web/`, `cloud/rdf/`, generators |

**Recommended order** (cheapest / highest value first):

1. **`boltdb/bolt` → `go.etcd.io/bbolt`** — bbolt is an API-compatible fork, so this is largely an import-path change across 2 files. Highest value: it is the local persistence layer, and the original has been archived for 7 years with bug fixes landing only in the fork.
2. **`yaml.v2` → `go.yaml.in/yaml/v3`** — 1 file, and the module is already in the tree.
3. **`oklog/ulid` → `/v2`** — 2 files, straightforward.
4. **`ioprogress`** — 1 file; consider inlining rather than depending on a 2018 orphan.
5. **`src-d/go-git.v4` → `go-git/go-git/v5`** — largest effort (the API changed between v4 and v5) but by far the biggest transitive payoff (see D6a). Target v5; v6 exists but is still alpha (`v6.0.0-alpha.5`).
6. **`triplestore`** — no migration path exists. Decide whether to fork it under this org or vendor it, since it is a core dependency with no upstream maintainer.

---

### D6a: `go-git` v4 drags in 26 stale transitive dependencies — **FIXED**

**Severity:** Medium  
**File:** `go.mod` (indirect block)

`gopkg.in/src-d/go-git.v4` alone pulls in 26 transitive modules, nearly all from 2015–2019:

```
github.com/alcortesm/tgz@v0.0.0-20161220082320     github.com/pkg/errors@v0.8.1
github.com/anmitsu/go-shlex@v0.0.0-20161002113705  github.com/sergi/go-diff@v1.0.0
github.com/armon/go-socks5@v0.0.0-20160902184237   github.com/src-d/gcfg@v1.4.0
github.com/emirpasic/gods@v1.12.0                  github.com/stretchr/objx@v0.2.0
github.com/flynn/go-shlex@v0.0.0-20150515145356    github.com/xanzy/ssh-agent@v0.2.1
github.com/gliderlabs/ssh@v0.2.2                   golang.org/x/crypto@v0.0.0-20190701094942
github.com/jessevdk/go-flags@v1.4.0                golang.org/x/net@v0.0.0-20190724013045
github.com/kevinburke/ssh_config@v0.0.0-20190725   gopkg.in/src-d/go-billy.v4@v4.3.2
```

Notable: it requests `golang.org/x/crypto@v0.0.0-20190701` and `golang.org/x/net@v0.0.0-20190724`. Those are only overridden to current versions because this project requires them **directly** at `v0.48.0` / `v0.49.0`. If those direct requirements were ever dropped or relaxed, MVS would silently select 2019-era crypto.

It also transitively pulls `github.com/pkg/errors@v0.8.1`, itself archived and superseded by stdlib `%w` wrapping (Go 1.13). Confirmed **not** used directly by this project — transitive only.

**Fix:** Migrating to `go-git/go-git/v5` (see D6) collapses most of this tree, since v5 uses maintained forks (`go-git/gcfg`, `go-git/go-billy/v5`) and current `x/crypto`. This is the single highest-leverage dependency change available.

**Verification after migration:**

```sh
go mod graph | grep -c "src-d"   # expect 0
go mod tidy && git diff go.mod go.sum
```

---

### D6b: Classic ELB (v1) support maintained alongside ELBv2

**Severity:** Low  
**Files:** `aws/spec/classicloadbalancer.go`, `aws/fetch/config.go`, `gen/aws/fetchers_definitions.go`

The project depends on both `elasticloadbalancing` (Classic ELB, v1) and `elasticloadbalancingv2` (ALB/NLB). Classic Load Balancers are a legacy AWS product, and EC2-Classic itself was fully retired in Aug 2022.

This is not broken and not urgent — Classic ELBs still exist in older VPC accounts, so the code has legitimate users. Flagged only so the ongoing maintenance cost is a conscious choice rather than inertia.

**Fix:** No action required. If Classic ELB support is ever dropped, it removes one SDK module plus the `classicloadbalancer.go` spec.

---

### D7: `release.go` references upstream wallix/awless paths — **FIXED** (release.go deleted, I11)

**Severity:** Low  
**File:** `release.go:105`

Build-info ldflags inject `github.com/wallix/awless/config` paths, and the upgrade URL mentions the wallix repo. This is technically correct for the module path but confusing since the fork is at babyhuey.

**Fix:** Audit all user-facing strings that reference `wallix/awless` and decide whether to update them (requires a module path change) or add fork-specific messaging.

---

### D8: No integration test coverage for SDK v2 API calls

**Severity:** Medium  
**Files:** `aws/fetch/manual_fetchers.go` (38KB), `aws/spec/*.go` (command implementations)

The only test coverage for AWS API interactions is via the broken acceptance framework (D1). Unit tests in `aws/conv/` test conversion logic, but actual API call orchestration (pagination, error handling, retries) is untested.

**Fix:** Introduce interface-based mocking (see D1) or use `smithy-go/middleware` test helpers to verify API call patterns.

---

### D9: Skipped tests

**Severity:** Low  
**Files:**
- `console/displayer_test.go:136` — skipped with comment about ordering
- `aws/spec/doc_test.go:11` — unconditionally skipped with `t.Skip()`

Tests that are permanently skipped are dead weight and may mask regressions.

**Fix:** Fix the underlying issues and unskip, or delete if no longer relevant.

---

## Improvements

### I19: Findings deferred by the golangci-lint v2 migration

**Severity:** Low (one exception, noted below)  
**File:** `.golangci.yml`

golangci-lint v2 merged `gosimple` and `stylecheck` into `staticcheck`. The v1 config enabled `gosimple` + `staticcheck` but **not** `stylecheck`, so `ST*` checks had never run against this codebase. Separately, `govet`'s `enable-all: true` now picks up a recently added `inline` analyzer that is not part of stock `go vet`.

Enabling all of it produces **271 findings**, so `B5` excluded them to keep that commit a behavior-preserving migration rather than a 271-fix changeset. Exclusions are `-ST*`, `-QF*`, and `govet: disable: [inline]`, each commented in place.

| Check | Count | What it is |
|---|---|---|
| `ST1003` | 175 | Naming — mostly `ALL_CAPS` constants (e.g. `template/env/env.go`) |
| `QF1012` | 25 | `Writer.Write([]byte(fmt.Sprintf(...)))` → `fmt.Fprintf` |
| `ST1005` | 20 | Error strings capitalised or ending in punctuation |
| `QF1008` | 19 | Redundant embedded-field selectors |
| `inline` (govet) | 16 | `reflect.Ptr` → `reflect.Pointer` (deprecated alias) |
| `ST1016` | 7 | Inconsistent receiver names |
| `QF1004` | 5 | `strings.Replace(s, a, b, -1)` → `strings.ReplaceAll` |
| `ST1012` | 2 | Error variable naming (`errFoo`) |
| `QF1002` | 1 | Untagged switch could be tagged |
| **`QF1009`** | **1** | **`time.Time` compared with `==` instead of `.Equal()`** |

**`QF1009` is the one worth attention** — `aws/tailers/scaling_activity.go:102` compares `time.Time` values with `==`, which compares wall clock, monotonic reading, *and* location. Two instants that represent the same moment can compare unequal. Worth fixing on its own merits rather than waiting for a lint sweep.

**Fix:** address these as part of `I9`'s incremental linter expansion, after the `%w` sweep in `I8`. Order by ratio of value to churn: `inline` (16, mechanical, removes deprecated API use) → `QF1009` (1, correctness) → `QF1004`/`QF1008`/`QF1012` (49, mechanical) → `ST1005`/`ST1012`/`ST1016` (29) → `ST1003` (175, largest and purely cosmetic; consider leaving `ALL_CAPS` exported identifiers alone since renaming them is a breaking API change).

Note `ST1003` on exported names cannot be fixed without breaking the public API of `template/env`, so that check may warrant staying disabled permanently rather than being treated as debt.

---

### I1: Add `arm64` to release builds — **FIXED** (GoReleaser, I11)

**Severity:** Low  
**File:** `release.go`

The `builds` map only includes `amd64` and `386`. The CI workflow already cross-compiles for `arm64`, but the release script doesn't produce `arm64` artifacts.

**Fix:** Add `"arm64"` to the `darwin` and `linux` build targets in `release.go`.

---

### I2: Add golangci-lint to CI via official action — **WON'T FIX**: pinned `go install` of an exact version (v2.12.2) is equivalent and keeps CI and local in sync; see B5

**Severity:** Low  
**File:** `.github/workflows/ci.yml`

The lint job builds golangci-lint from source (`go install ... @latest`). This is slow and unpinned. The official `golangci/golangci-lint-action` is faster (uses pre-built binaries) and supports version pinning.

**Fix:** Replace with `golangci/golangci-lint-action@v6` with a pinned version.

---

### I3: Add CRUD commands for new services

**Severity:** Medium  
**Scope:** 8 services × common operations

Extend EKS, DynamoDB, Secrets Manager, API Gateway v2, SSM, EFS, CloudTrail, and CloudWatch Logs with create/update/delete commands. Priority by user impact:

1. **Secrets Manager:** `create secret`, `update secret`, `delete secret`
2. **SSM:** `create ssmparameter`, `update ssmparameter`, `delete ssmparameter`
3. **EKS:** `create ekscluster`, `delete ekscluster`, `create eksnodegroup`, `delete eksnodegroup`
4. **DynamoDB:** `create dynamodbtable`, `delete dynamodbtable`
5. **EFS:** `create filesystem`, `delete filesystem`

---

### I4: Implement proper signal handling for graceful shutdown

**Severity:** Low  
**File:** `main.go`, `commands/root.go`

No `SIGINT`/`SIGTERM` handling. If a long-running AWS operation (e.g., `create instance` waiting for status) is interrupted, there's no cleanup or cancellation. Combined with D4 (no context propagation), Ctrl+C just kills the process.

**Fix:** Install signal handlers in `main.go`, create a root context, and propagate it through the command chain.

---

### I5: Modernize the code generation pipeline

**Severity:** Low  
**Files:** `gen/aws/generators/*.go`

The generators use raw `text/template` with string concatenation and manual file writing. Modern Go would use `go generate` directives more idiomatically, and the templates could benefit from:
- Better error messages on template failures
- Validation that generated output compiles before overwriting
- A `--check` mode that verifies generated files are up-to-date (useful in CI)

---

### I6: Add Dependabot auto-merge for patch updates

**Severity:** Low  
**File:** `.github/dependabot.yml`

Dependabot is configured but there's no auto-merge setup. Given the high number of AWS SDK dependencies (30+ service modules), patch updates create significant PR noise.

**Fix:** Add a GitHub Actions workflow that auto-merges Dependabot PRs for patch-level updates after CI passes.

---

### I7: Test coverage improvement targets

**Severity:** Medium  
**Current state:** ~22.4% coverage (83 test files, 180 source files)

Key untested areas:
- `commands/` — only `tabcompletion_test.go` and `run_test.go` exist; no tests for `ssh.go`, `show.go`, `list.go`, `sync.go`
- `aws/fetch/manual_fetchers.go` (38KB) — zero test coverage
- `aws/services/` — only lookup and relation tests; no service-level tests
- `ssh/` — basic tests exist but no integration-style tests

**Fix:** Prioritize `commands/` tests (testable with mock services) and `aws/conv/` edge cases.

---

### I8: Go language modernization — **FIXED**

**Severity:** Medium  
**Scope:** Project-wide

The codebase was originally written for Go 1.9 and while it compiles on Go 1.26, it doesn't leverage language and stdlib improvements from the last 8 major releases:

**`interface{}` → `any` (Go 1.18)**
- 527 occurrences of `interface{}` in non-generated, non-test code
- Direct alias replacement, zero behavior change

**Error wrapping with `%w` (Go 1.13)**
- 0 uses of `fmt.Errorf("...: %w", err)` anywhere in the codebase
- 127 uses of `fmt.Errorf("...: %s", err)` which discard the error chain
- This breaks `errors.Is()` and `errors.As()` for callers — error inspection is impossible through these wrapping points

**`slices` and `maps` packages (Go 1.21)**
- 10+ uses of `sort.Slice` / `sort.Strings` that could use `slices.SortFunc` / `slices.Sort`
- Index-based loops (`for i := 0; i < len(x); i++`) where `slices.Contains`, `slices.Index`, or range would suffice

**Range over integers (Go 1.22)**
- Several `for i := 0; i < n; i++` loops iterating over `reflect.Value.Len()` or slice lengths that could use `range n` syntax

**Structured logging with `log/slog` (Go 1.21)**
- Custom `logger/` package with ad-hoc formatting
- Could adopt `slog` for structured output while keeping the existing UX for interactive mode

**`smithy.APIError` type assertions → `errors.As` (Go 1.13)**
- 8+ manual type assertions (`err.(smithy.APIError)`) in `aws/spec/*.go`
- Should use `errors.As(&apiErr)` which handles wrapped errors correctly

**Priority order for migration:**
1. `%s` → `%w` for error wrapping (fixes a real semantic bug — error chains are silently broken)
2. `interface{}` → `any` (mechanical, high-value readability improvement)
3. `smithy.APIError` assertions → `errors.As` (correctness)
4. `sort.*` → `slices.*` (minor readability improvement)
5. `log/slog` adoption (larger effort, optional)

---

### I9: Expand golangci-lint configuration — **FIXED**

**Severity:** Medium  
**File:** `.golangci.yml`

The current config enables only 10 linters. Several high-value linters are missing that would catch real bugs:

| Linter | What it catches | Estimated hits |
|--------|----------------|----------------|
| `errcheck` | Unchecked error returns (e.g., `file.Close()`, `writer.Flush()`) | ~10+ in ssh, database, console |
| `errorlint` | Non-wrapping `%s` in `fmt.Errorf`, type assertions on errors instead of `errors.As` | 127+ (see I8) |
| `noctx` | HTTP requests without `context.Context` (`http.Get`, `http.NewRequest` without context) | 2 (`commands/run.go`, `config/upgrade.go`) |
| `bodyclose` | HTTP response bodies not closed | Likely clean (checked manually) |
| `nilerr` | Returning `nil` when err is non-nil (silent error swallowing) | Unknown |
| `gocritic` | Broad set of opinionated but high-signal checks (dupSubExpr, sloppyLen, etc.) | Unknown |
| `prealloc` | Slice allocations that could be pre-sized | Low priority, perf-only |
| `exhaustive` | Non-exhaustive switch statements on enums | Useful for cloud resource types |
| `revive` | Superset of `golint` with configurable rules | Many, needs tuning |
| `wrapcheck` | Errors from external packages returned without wrapping | Very noisy initially |

**Recommended rollout order:**
1. `errcheck` — catches real resource leaks (unclosed files, unflushed writers)
2. `errorlint` — enforces `%w` and `errors.As` (complements I8)
3. `noctx` — trivial to fix (only 2 hits), catches context-less HTTP calls
4. `gocritic` — broad quality improvement with low false-positive rate
5. `exhaustive` — catches missing cases in type switches on AWS resource types

**Note:** Enable incrementally. Adding all at once would produce hundreds of findings and block CI. Use `issues.new-from-rev` or `issues.max-issues-per-linter` to ratchet.

---

### I10: Add govulncheck to CI — **FIXED**

**Severity:** High  
**File:** `.github/workflows/ci.yml`

No vulnerability scanning exists. The project has 60+ dependencies including security-sensitive ones:

- `golang.org/x/crypto v0.48.0` — used for SSH
- `gopkg.in/src-d/go-git.v4 v4.13.1` — **abandoned**, last release 2019, known CVEs in the go-git ecosystem
- `github.com/boltdb/bolt v1.3.1` — **unmaintained** (replaced by `go.etcd.io/bbolt`)
- 30+ AWS SDK v2 modules

**What govulncheck provides:**
- Checks Go module dependencies against the official [Go vulnerability database](https://vuln.go.dev)
- Reports only vulnerabilities in code paths actually used (not just imported)
- Zero config required — just `govulncheck ./...`

**Recommended CI addition:**

```yaml
  vuln:
    name: Vulnerability Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - name: Run govulncheck
        run: govulncheck ./...
```

**Known high-risk dependencies:**
- `gopkg.in/src-d/go-git.v4` — should migrate to `github.com/go-git/go-git/v5` regardless of vuln findings (see D6)
- `github.com/boltdb/bolt` — should migrate to `go.etcd.io/bbolt` (maintained fork, API-compatible)

---

### I11: Replace `release.go` with GoReleaser — **FIXED**

**Severity:** Medium  
**File:** `release.go`

The current release process is a 200-line hand-rolled Go script (`release.go`) that:
- Only builds `amd64` and `386` (no `arm64`)
- Manually creates `.tar.gz` and `.zip` archives
- Has no checksum file generation
- Doesn't create GitHub releases or upload artifacts
- Has stale Homebrew references to the wallix upstream
- Uses `go build` without CGO_ENABLED=0 (non-static binaries)
- Has no changelog generation
- Requires manual invocation (`go run release.go -tag v1.0.0`)

**GoReleaser replaces all of this** with a single `.goreleaser.yml` and a GitHub Actions workflow:

```yaml
# .goreleaser.yml
version: 2

project_name: awless

before:
  hooks:
    - go mod tidy

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/wallix/awless/config.buildDate={{.Date}}
      - -X github.com/wallix/awless/config.buildSha={{.Commit}}
      - -X github.com/wallix/awless/config.buildOS={{.Os}}
      - -X github.com/wallix/awless/config.buildArch={{.Arch}}

archives:
  - formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats:
          - zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^ci:'

release:
  github:
    owner: babyhuey
    name: awless
```

**CI release workflow** (`.github/workflows/release.yml`):

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**What this gets you:**
- arm64 builds for macOS and Linux (Apple Silicon, AWS Graviton)
- Automatic GitHub Release with changelog from commit messages
- SHA256 checksums file
- Reproducible, static binaries (CGO_ENABLED=0)
- Homebrew tap support (add `brews:` section)
- Docker image publishing (add `dockers:` section)
- Snapcraft/scoop/AUR support if needed later
- `release.go` can be deleted entirely

**Migration effort:** Low. Drop-in replacement — just add the two files above and delete `release.go`.

---

### I12: Restore `-race` in CI — **FIXED**

**Severity:** Medium  
**File:** `.github/workflows/ci.yml`

The race detector was deliberately removed to get CI green (commit `6830d3d0` — "drop -race flag"). History also shows a prior data race that required a dependency bump to fix (`738b5e0a` — "Vendoring latest triplestore to fix datarace issue").

This is the worst combination: a codebase with heavy goroutine fan-out (20+ `go func` sites across `sync/`, `fetch/`, `aws/fetch/`, `aws/conv/`), a documented history of data races, and no race detection in CI.

**Fix:** Add a separate CI job running `go test -race ./...`. Keep it non-blocking (`continue-on-error: true`) initially so it surfaces races without gating merges, then make it required once the existing findings are fixed. Running it as its own job avoids slowing down the main test job.

```yaml
  race:
    name: Race Detector
    runs-on: ubuntu-latest
    continue-on-error: true   # remove once existing races are fixed
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - run: go test -race -count=1 ./...
```

---

### I13: Unbounded concurrency in fan-out fetchers — **FIXED**

**Severity:** Medium  
**Files:** same as `B8`; both were fixed by the same change.

Unbounded fan-out meant an account with thousands of buckets, clusters, IAM
users, queues or hosted zones issued that many simultaneous AWS API calls,
reliably causing throttling (`RequestLimitExceeded`) and unpredictable memory
use.

Bounded via `errgroup.SetLimit`. `aws/fetch/limits.go` defines
`maxParallelAWSCalls = 20`; `inspect/inspectors/pricer.go` uses a local limit of
10 for the pricing endpoint.

**Possible follow-up:** make the limit configurable through `awless config` if 20
proves wrong for very large accounts.

---

### I14: Migrate fuzzers to native Go fuzzing

**Severity:** Low  
**Files:** `template/fuzz/parsing/main.go:9`, `template/fuzz/parameters/main.go:12`

Both fuzzers use the legacy `go-fuzz` signature:

```go
func Fuzz(data []byte) int
```

Native fuzzing has been in the toolchain since Go 1.18 and uses `func FuzzXxx(f *testing.F)`. The legacy form requires the external `go-fuzz` tool and cannot run under `go test`, so these fuzzers almost certainly never execute — they aren't wired into CI and the existing corpus in `template/fuzz/*/workdir/corpus/` sits unused.

**Fix:** Convert to native fuzz targets in the packages under test, seed them from the existing corpus files, and add a short fuzz run to CI:

```go
func FuzzParseTemplate(f *testing.F) {
    f.Add([]byte("create instance type=t2.micro"))
    f.Fuzz(func(t *testing.T, data []byte) {
        _, _ = template.Parse(string(data))  // must not panic
    })
}
```

```yaml
      - run: go test -run=XXX -fuzz=FuzzParseTemplate -fuzztime=60s ./template/...
```

The template parser is generated from a PEG grammar and takes untrusted file input via `awless run`, so it's a genuinely good fuzz target. Note the parser already has a `recover()` in `template/parser.go:33`, which suggests panics have been hit before.

---

### I15: Generated code is out of sync and the generators produce non-compiling output

**Severity:** Medium — raised from Low; both problems are now measured, not hypothetical  
**Files:** `gen/aws/generators/*.go`, all `gen_*.go`, `.github/workflows/ci.yml`

Running the generators against the current tree revealed two concrete problems.

**1. The generators' output does not compile.**

```
$ cd gen/aws/generators && go run *.go
$ go vet ./aws/services/
vet: aws/services/gen_mocks_test.go:39:3: "github.com/aws/aws-sdk-go-v2/service/s3" imported and not used
```

The mock template emits an import it does not use, so a fresh regeneration
breaks the build. The committed `gen_mocks_test.go` compiles, meaning it was
either produced by an older generator or hand-corrected — either way the
generators are not currently a usable source of truth.

**2. ~2,400 lines of pre-existing drift.**

Regenerating at an untouched checkout produces `3 files changed, 2381
insertions(+), 1946 deletions(-)`, almost all in `aws/spec/gen_runs.go`. So the
committed generated files have not matched their generators for some time.

This is why the module rename in `78e9316f` rewrote generated files in place
instead of regenerating: regenerating would have mixed an unrelated 2,400-line
drift correction, plus a compile failure, into that change.

**Fix, in order:**

1. Fix the unused-import bug in the mock generator template so output compiles.
2. Regenerate everything and commit the drift correction as its own change, so
   the diff is reviewable in isolation.
3. Only then add the CI drift check below, which is meaningless until 1 and 2
   land.

```yaml
  codegen:
    name: Codegen Drift
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go install golang.org/x/tools/cmd/goimports@latest
      - run: cd gen/aws/generators && go run *.go
      - run: git diff --exit-code || (echo "Generated files are out of date. Run: cd gen/aws/generators && go run *.go" && exit 1)
```

This also enforces the "never hand-edit `gen_*.go`" rule documented in
`AGENTS.md`, which the rename commit had to bend precisely because the
generators were unreliable.

---

### I16: Replace `os.Exit` calls in command paths

**Severity:** Low  
**Files:** `commands/errors.go:29`, `commands/runner.go:89`, `commands/show.go:77,260`, `commands/list.go:119,123`, `commands/revert.go:60`, `commands/run.go:199`, `template/runner.go:93`, `aws/config/validator.go:63`

Nine `os.Exit` calls sit inside command implementations. `os.Exit` skips all deferred functions, which means:

- BoltDB is not closed cleanly (`database/db.go:44` uses `defer db.Close()`)
- Open file handles and SSH sessions are not released
- Template execution logs may not be flushed

It also makes the affected code paths untestable, which is part of why `commands/` has almost no test coverage (see I7).

**Fix:** Return errors up to `main.go` and exit once from there. `cobra` supports `RunE` for exactly this. This pairs with I4 (signal handling), since both require a single well-defined exit path.

---

### I17: `main.go` discards the error from `Execute()` — **FIXED**

**Severity:** Low  
**File:** `main.go:22`

```go
func main() {
    commands.RootCmd.Execute()
}
```

The returned error is dropped, so the process exits `0` even when a command fails at the cobra level. This breaks shell scripting, CI pipelines, and `&&` chaining around `awless`. The codebase compensates with scattered `os.Exit(1)` calls (see I16), but any failure path that returns an error instead of calling `os.Exit` silently reports success.

**Fix:**

```go
func main() {
    if err := commands.RootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

Note `errcheck` (proposed in I9) would have caught this.

---

### I18: Pin the Go toolchain version — **PARTIALLY DONE**

**Severity:** Low  
**Files:** `go.mod:3`, `.github/workflows/ci.yml`

**Done:** all three CI jobs now use `go-version-file: go.mod` instead of a hardcoded `'1.26'`, so the Go version has a single source of truth and CI installs exactly `1.26.1` rather than a floating latest patch. The `test` job's single-value `go-version` matrix was removed as redundant; the `build` job's `goos`/`goarch` matrix is untouched.

**Deliberately not done:** no `toolchain` directive was added to `go.mod`. It would pin local developer builds to an exact patch, but the dev toolchain in use here is a custom build (`go1.26.5-X:nodwarf5`), and a `toolchain` line risks Go fetching a different official toolchain in preference to it. CI is already deterministic without it.

**Remaining (optional):** add `toolchain go1.26.x` if reproducible *release* artifacts become a requirement — most relevant alongside `I11` (GoReleaser), where the build toolchain forms part of an artifact's provenance.
