# Local development commands.
#
# Tool versions are pinned here and in .github/workflows/ci.yml; keep them in
# sync. Run `make help` to list targets, `make verify` for the full CI gate.

SHELL := /bin/bash
.DEFAULT_GOAL := help

BINARY          := awless
MODULE          := github.com/bootswithdefer/awless
GOBIN           := $(shell go env GOPATH)/bin

# Must match the version installed by the lint job in CI.
GOLANGCI_VERSION := v2.12.2
GOLANGCI         := $(GOBIN)/golangci-lint
GORELEASER       := $(GOBIN)/goreleaser
GOVULNCHECK      := $(GOBIN)/govulncheck
GOIMPORTS        := $(GOBIN)/goimports

# Files that are generated or vendored grammar output, excluded from formatting
# checks to match the exclusions in .golangci.yml.
FMT_FILES := $(shell git ls-files '*.go' | grep -v 'gen_.*\.go$$' | grep -v '\.peg\.go$$')

.PHONY: help
help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

## --- build ---------------------------------------------------------------

.PHONY: build
build: ## Build the awless binary
	go build -o $(BINARY) .

.PHONY: install
install: ## Install awless into GOPATH/bin
	go install .

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out

## --- test ----------------------------------------------------------------

.PHONY: test
test: ## Run tests
	go test -count=1 ./...

.PHONY: test-race
test-race: ## Run tests with the race detector
	go test -race -count=1 ./...

.PHONY: fuzz
fuzz: ## Run each fuzz target briefly (FUZZTIME=30s to override)
	go test ./template/ -run=XXX -fuzz='^FuzzParse$$' -fuzztime=$(or $(FUZZTIME),30s)
	go test ./template/ -run=XXX -fuzz='^FuzzParseStatement$$' -fuzztime=$(or $(FUZZTIME),30s)

.PHONY: cover
cover: ## Run tests with coverage and print the total
	go test -count=1 -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

.PHONY: cover-html
cover-html: cover ## Open the HTML coverage report
	go tool cover -html=coverage.out

## --- static analysis -----------------------------------------------------

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: $(GOLANGCI) ## Run golangci-lint
	$(GOLANGCI) run ./...

.PHONY: lint-config
lint-config: $(GOLANGCI) ## Validate .golangci.yml against its schema
	$(GOLANGCI) config verify

.PHONY: vuln
vuln: $(GOVULNCHECK) ## Check for known vulnerabilities reachable from this code
	$(GOVULNCHECK) ./...

## --- formatting ----------------------------------------------------------

.PHONY: fmt
fmt: $(GOIMPORTS) ## Format code (gofmt -s + goimports with local prefix)
	gofmt -w -s $(FMT_FILES)
	$(GOIMPORTS) -w -local $(MODULE) $(FMT_FILES)

.PHONY: fmt-check
fmt-check: ## Fail if any non-generated file is unformatted
	@unformatted=$$(gofmt -l -s $(FMT_FILES)); \
	if [ -n "$$unformatted" ]; then \
		echo "Not gofmt'd:"; echo "$$unformatted"; exit 1; \
	fi
	@echo "all files formatted"

## --- codegen -------------------------------------------------------------

.PHONY: generate
generate: ## Regenerate gen_*.go from gen/aws definitions and templates
	cd gen/aws/generators && go run *.go

.PHONY: generate-check
generate-check: generate ## Fail if regeneration produces a diff (codegen drift)
	@git diff --exit-code || { \
		echo "Generated files are out of date. Commit the regeneration."; exit 1; }

## --- aggregate gates -----------------------------------------------------

.PHONY: check
check: fmt-check vet lint test ## Fast pre-commit gate

.PHONY: verify
verify: fmt-check vet lint test-race vuln ## Full gate, mirrors CI
	@echo
	@echo "verify: all checks passed"

## --- tooling -------------------------------------------------------------

.PHONY: tools
tools: $(GOLANGCI) $(GOVULNCHECK) $(GOIMPORTS) ## Install pinned dev tools

$(GOLANGCI):
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

$(GOVULNCHECK):
	go install golang.org/x/vuln/cmd/govulncheck@latest

$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: release-check
release-check: $(GORELEASER) ## Validate .goreleaser.yml
	$(GORELEASER) check

.PHONY: release-snapshot
release-snapshot: $(GORELEASER) ## Build release artifacts locally without publishing
	$(GORELEASER) release --snapshot --clean --skip=publish

$(GORELEASER):
	go install github.com/goreleaser/goreleaser/v2@latest

.PHONY: pinact
pinact: ## Re-pin GitHub Actions to commit SHAs (requires pinact)
	pinact run

.PHONY: pinact-update
pinact-update: ## Update GitHub Actions to latest versions and re-pin
	pinact run --update

.PHONY: tidy
tidy: ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: hooks
hooks: ## Enable the repo's pre-commit hooks
	git config core.hooksPath .githooks
