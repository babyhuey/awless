//go:generate go run $GOFILE properties.go mocks.go fetchers.go services.go commands.go acceptance_mocks.go

//go:generate gofmt -s -w ../../../aws
//go:generate goimports -w ../../../aws

//go:generate gofmt -s -w ../../../aws/services
//go:generate goimports -w ../../../aws/services

//go:generate gofmt -s -w ../../../aws/fetch
//go:generate goimports -w ../../../aws/fetch

//go:generate gofmt -s -w ../../../cloud/properties
//go:generate goimports -w ../../../cloud/properties

//go:generate gofmt -s -w ../../../cloud/rdf
//go:generate goimports -w ../../../cloud/rdf

//go:generate gofmt -s -w ../../../aws/spec
//go:generate goimports -w ../../../aws/spec

//go:generate gofmt -s -w ../../../acceptance/aws
//go:generate goimports -w ../../../acceptance/aws

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"text/template"
)

// localPrefix matches the goimports -local setting in .golangci.yml and the
// Makefile, so generated files group imports the same way as hand-written ones.
const localPrefix = "github.com/bootswithdefer/awless"

var (
	ROOT_DIR = filepath.Join("..", "..", "..")

	FETCHERS_DIR         = filepath.Join(ROOT_DIR, "aws", "fetch")
	SERVICES_DIR         = filepath.Join(ROOT_DIR, "aws", "services")
	SPEC_DIR             = filepath.Join(ROOT_DIR, "aws", "spec")
	AWSAT_DIR            = filepath.Join(ROOT_DIR, "acceptance", "aws")
	CLOUD_PROPERTIES_DIR = filepath.Join(ROOT_DIR, "cloud", "properties")
	CLOUD_RDF_DIR        = filepath.Join(ROOT_DIR, "cloud", "rdf")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[+] ")
	flag.Parse()

	// fetchers
	generateFetcherFuncs()
	generateServicesFuncs()

	// mocks
	generateTestMocks()

	// commands
	generateCommands()
	generateAcceptanceMocks()
	generateAcceptanceFactory()

	// properties
	generateProperties()
	generateRDFProperties()
}

func writeTemplateToFile(templ *template.Template, data any, dir, filename string) {
	path := filepath.Join(dir, filename)

	var buff bytes.Buffer
	if err := templ.Execute(&buff, data); err != nil {
		// Name the target file: a bare template error gives no clue which of the
		// ten generated files failed.
		log.Fatalf("executing template for %s: %s", relativePathToRoot(path), err)
	}

	// Parse before overwriting. A template change that produces invalid Go used to
	// land on disk and only fail at the next build, by which point the previous
	// good copy was gone.
	if err := validateGoSyntax(path, buff.Bytes()); err != nil {
		log.Fatalf("generated output for %s is not valid Go: %s", relativePathToRoot(path), err)
	}

	if err := os.WriteFile(path, buff.Bytes(), 0644); err != nil {
		log.Fatalf("writing %s: %s", relativePathToRoot(path), err)
	}

	// Templates emit an import for every AWS API in the definitions, whether or
	// not the generated body references it — the s3 client package, for example,
	// is unused when only s3types is referenced. That produced output that did
	// not compile.
	//
	// goimports drops the unused imports and groups the rest. It runs here rather
	// than only via the //go:generate directives above, so `go run *.go` produces
	// compiling output too.
	if err := runGoimports(path); err != nil {
		log.Fatalf("goimports %s: %s", relativePathToRoot(path), err)
	}

	log.Printf("generated %s", relativePathToRoot(path))
}

// validateGoSyntax parses src, so a template producing invalid Go fails before
// the previous good file is overwritten. It checks syntax only; type errors still
// surface at build time.
func validateGoSyntax(path string, src []byte) error {
	_, err := parser.ParseFile(token.NewFileSet(), path, src, parser.SkipObjectResolution)
	return err
}

// runGoimports formats a generated file and prunes its unused imports.
func runGoimports(path string) error {
	bin, err := exec.LookPath("goimports")
	if err != nil {
		return fmt.Errorf("goimports not found on PATH; install it with `go install golang.org/x/tools/cmd/goimports@latest`: %w", err)
	}
	out, err := exec.Command(bin, "-w", "-local", localPrefix, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}
	return nil
}

func relativePathToRoot(path string) string {
	rel, _ := filepath.Rel(ROOT_DIR, path)
	return rel
}

// capitalize upper-cases the first character of s.
//
// Replaces strings.Title, deprecated in Go 1.18 because it applies Unicode word
// boundaries and title-cases every word. Every input here is a single ASCII
// token — an AWS API name, resource type, template action, or policy effect —
// so this is both correct and narrower than the deprecated behavior.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
