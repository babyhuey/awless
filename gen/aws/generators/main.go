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
	"log"
	"os"
	"path/filepath"
	"strings"

	"text/template"
)

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
	var buff bytes.Buffer
	if err := templ.Execute(&buff, data); err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, buff.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	log.Printf("generated %s", relativePathToRoot(path))
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
