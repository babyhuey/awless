# AGENTS.md

Guide for AI agents working in this repository.

## Project Overview

`awless` is a CLI tool for managing AWS resources. It provides a template DSL for infrastructure creation/revert, local graph-based resource sync, smart SSH, and human-friendly output. This is a modernized fork of [wallix/awless](https://github.com/wallix/awless) maintained at [bootswithdefer/awless](https://github.com/bootswithdefer/awless).

- **Language:** Go 1.26
- **Module path:** `github.com/bootswithdefer/awless`
- **AWS SDK:** v2 (`github.com/aws/aws-sdk-go-v2`)
- **CLI framework:** `github.com/spf13/cobra`
- **Graph storage:** `github.com/bootswithdefer/triplestore` (RDF triples)
- **Local DB:** bbolt (`go.etcd.io/bbolt`)
- **Version:** v1.0.0 (set in `config/version.go`)

## Build & Test

```sh
# Build
go build -o awless .

# Run tests
go test ./...

# Lint (requires golangci-lint)
golangci-lint run ./...

# Format
gofmt -w -s .
goimports -w -local github.com/bootswithdefer/awless .

# Code generation (regenerates gen_*.go files)
cd gen/aws/generators && go run *.go

# Full check
make check
```

## Directory Structure

```
.
├── main.go                  # Entry point — delegates to commands.RootCmd.Execute()
├── commands/                # Cobra CLI commands (list, run, show, ssh, sync, etc.)
├── aws/
│   ├── services/            # Service implementations (gen_services.go is generated)
│   ├── spec/                # Command specs per resource (gen_runs.go, gen_cmds_defs.go, gen_inits.go generated)
│   ├── fetch/               # Data fetchers (gen_fetchers.go is generated)
│   ├── conv/                # AWS SDK type → internal model conversion
│   ├── config/              # AWS config validation
│   ├── doc/                 # CLI documentation helpers
│   └── tailers/             # CloudFormation/ASG event tailers
├── cloud/                   # Cloud abstraction layer (interfaces, properties, RDF)
│   ├── properties/          # Generated property constants
│   └── rdf/                 # Generated RDF namespace constants
├── template/                # Template DSL engine
│   ├── internal/ast/        # PEG grammar + parser (awless-template-syntax.peg)
│   ├── env/                 # Template execution environment
│   ├── params/              # Parameter validation
│   └── fuzz/                # Fuzz testing corpus
├── graph/                   # RDF-based resource graph
├── gen/aws/                 # Code generation
│   ├── generators/          # Generator programs (go run *.go)
│   ├── properties_definitions.go
│   ├── fetchers_definitions.go
│   └── mock_definitions.go
├── acceptance/aws/          # Acceptance test framework (currently broken — SDK v2 mocking TODO)
├── config/                  # App config, versioning, upgrade logic
├── console/                 # Terminal display, table formatting, column headers
├── database/                # BoltDB-backed local storage
├── fetch/                   # Generic fetch framework
├── inspect/                 # Infrastructure analysis inspectors
├── logger/                  # Custom logger
├── ssh/                     # SSH client implementation
├── sync/                    # Cloud → local graph sync
├── web/                     # Web-based resource viewer
├── smoke_tests/             # Shell-based integration tests (require AWS credentials)
├── .github/workflows/ci.yml # GitHub Actions CI
└── .githooks/pre-commit     # gofmt + golangci-lint pre-commit hook
```

## Code Generation

This project uses Go code generation extensively. Generated files follow the `gen_*.go` naming convention and are excluded from linting (see `.golangci.yml`).

**How to regenerate:**

```sh
cd gen/aws/generators && go run *.go
```

**What gets generated:**

| Definition file | Output |
|----------------|--------|
| `gen/aws/properties_definitions.go` | `cloud/properties/gen_properties.go`, `cloud/rdf/gen_rdf.go` |
| `gen/aws/fetchers_definitions.go` | `aws/fetch/gen_fetchers.go` |
| `gen/aws/mock_definitions.go` | `aws/services/gen_mocks_test.go`, `acceptance/aws/gen_mocks.go`, `acceptance/aws/gen_factory.go` |
| Definitions in generators/*.go | `aws/services/gen_services.go`, `aws/spec/gen_runs.go`, `aws/spec/gen_cmds_defs.go`, `aws/spec/gen_inits.go` |

**Do NOT edit `gen_*.go` files directly.** Edit the definitions in `gen/aws/` or the generator templates in `gen/aws/generators/`, then regenerate.

## Adding a New AWS Service

Follow this pattern (see `SERVICES_TODO.md` for candidates):

1. Add SDK dependency: `go get github.com/aws/aws-sdk-go-v2/service/<servicename>`
2. Define resources in `gen/aws/properties_definitions.go` (properties, RDF types)
3. Define fetchers in `gen/aws/fetchers_definitions.go`
4. Add service struct in `gen/aws/generators/services.go`
5. Add mock definition in `gen/aws/mock_definitions.go`
6. Run `cd gen/aws/generators && go run *.go` to regenerate
7. Implement manual fetchers in `aws/fetch/manual_fetchers.go` for complex APIs
8. Add conversion logic in `aws/conv/model.go` and `aws/conv/convert.go`
9. Register the service in `aws/services/init.go`
10. Add resource-specific command specs in `aws/spec/<resource>.go`
11. Wire display columns in `console/defaults.go` and `console/headers.go`

## Template DSL

The template language uses a PEG grammar at `template/internal/ast/awless-template-syntax.peg`. The compiled parser is `awless-template-syntax.peg.go` (generated by [pointlander/peg](https://github.com/pointlander/peg)).

Template syntax: `ACTION ENTITY param=value [param=value ...]`

Examples:
```
create instance type=t2.micro subnet=@my-subnet name=web-server
attach securitygroup id={instance.SecurityGroups} instance=$instance_id
```

**Do NOT edit the `.peg.go` file.** Edit the `.peg` file and regenerate with `peg`.

## Conventions

- **Imports:** Group stdlib, then third-party, then `github.com/bootswithdefer/awless` (enforced by goimports with `-local`)
- **Error handling:** Wrap with context; use `fmt.Errorf` or `decorateAWSError()` for AWS errors
- **Pointer helpers:** Use `String()`, `StringValue()`, `Int64()`, `Bool()` from `aws/spec/spec.go`
- **Resource types:** Constants defined in `cloud/cloud.go` (e.g., `cloud.Instance`, `cloud.Vpc`)
- **Service naming:** Lowercase single word (e.g., `infra`, `access`, `eks`, `dynamodb`)
- **Generated file prefix:** `gen_` — never hand-edit these
- **Test file suffix:** `_test.go` for unit tests, `_extra_test.go` for integration-style tests in external test packages

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`):
- **test:** `go test -count=1 -coverprofile=coverage.out ./...` on Go 1.26
- **lint:** `go vet ./...` then `golangci-lint run ./...`
- **build:** Cross-compile for linux/darwin/windows × amd64/arm64

## Linting Configuration

`.golangci.yml` enables: govet, ineffassign, staticcheck, unused, gosimple, gofmt, goimports, misspell, unconvert.

Notable exclusions:
- `SA1019` (deprecated usage) is suppressed — too noisy during SDK migration
- `gen_*.go` and `awless-template-syntax.peg.go` are excluded from linting
- `fieldalignment` and `shadow` disabled for govet

## Key Interfaces

- `cloud.Service` — implemented per AWS service group (access, infra, storage, etc.)
- `cloud.Resource` — generic resource with properties, used throughout graph/display
- `command` (in `aws/spec/spec.go`) — template command with `ParamsSpec()`, `inject()`, `Run()`
- `BeforeRunner` / `AfterRunner` / `ResultExtractor` — lifecycle hooks on commands

## Things to Watch Out For

- **Two compiled binaries in git:** `./generators` (root) and `gen/aws/generators/generators` — these should not be committed but currently are
- **Acceptance tests are broken:** All 30 factory functions in `acceptance/aws/gen_factory.go` have `// TODO: SDK v2 mocking needs rework`
- **`strings.Title` is deprecated:** Used in 5 generator files — produces lint warnings with newer Go versions
- **`.travis.yml` is stale:** References Go 1.9–1.11, completely superseded by GitHub Actions
- **Template PEG regeneration:** Requires the `peg` tool (`go install github.com/pointlander/peg@latest`)
