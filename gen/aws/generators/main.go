//go:generate go run $GOFILE properties.go mocks.go fetchers.go services.go commands.go acceptance_mocks.go entities.go

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
	rootDir = filepath.Join("..", "..", "..")

	fetchersDir        = filepath.Join(rootDir, "aws", "fetch")
	servicesDir        = filepath.Join(rootDir, "aws", "services")
	specDir            = filepath.Join(rootDir, "aws", "spec")
	templateASTDir     = filepath.Join(rootDir, "template", "internal", "ast")
	awsatDir           = filepath.Join(rootDir, "acceptance", "aws")
	cloudPropertiesDir = filepath.Join(rootDir, "cloud", "properties")
	cloudRDFDir        = filepath.Join(rootDir, "cloud", "rdf")
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
	generateTemplateEntities()

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
	rel, _ := filepath.Rel(rootDir, path)
	return rel
}

// capitalize renders s as an exported Go identifier.
//
// Known initialisms are upper-cased whole, so the "dns" service generates DNS
// rather than Dns. Without this, generated identifiers trip staticcheck's ST1003,
// and the hand-written functions they must match are forced to spell an initialism
// the way Go style says not to.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	if up, ok := initialisms[strings.ToLower(s)]; ok {
		return up
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// initialisms is keyed by the lower-case service or API name. Only entries that
// actually occur need listing; anything absent falls back to simple capitalization.
//
// "api" is deliberately absent: capitalize also renders SDK type names, and
// apigatewayv2's type is Api, so upper-casing it here produces a reference that
// does not exist.
var initialisms = map[string]string{
	"acm": "ACM",
	"api": "API",
	"cdn": "CDN",
	"dns": "DNS",
	"ecr": "ECR",
	"ecs": "ECS",
	"eks": "EKS",
	"iam": "IAM",
	"kms": "KMS",
	"rds": "RDS",
	"sns": "SNS",
	"sqs": "SQS",
	"ssm": "SSM",
	"sts": "STS",
	"s3":  "S3",
	"efs": "EFS",
}
