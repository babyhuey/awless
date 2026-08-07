.PHONY: build test lint fmt generate clean check

build:
	@echo Building application binary
	@go build -o awless .

test:
	@echo Running tests
	@go test ./...

lint:
	@echo Running linters
	@golangci-lint run ./...

fmt:
	@echo Formatting code
	@goimports -w -local github.com/bootswithdefer/awless .
	@gofmt -w -s .

generate:
	@echo Generating commands code: runtime, doc, etc.
	@cd gen/aws/generators && go run *.go

clean:
	@rm -f awless

check: fmt lint test
