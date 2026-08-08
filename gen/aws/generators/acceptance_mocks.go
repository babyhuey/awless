package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"text/template"
)

func generateAcceptanceMocks() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, specDir, func(os.FileInfo) bool { return true }, 0)
	if err != nil {
		panic(err)
	}

	finder := &findStructs{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Walk(finder, f)
		}
	}

	usedApis := make(map[string]bool)
	for _, cmd := range finder.result {
		if cmd.API == "" {
			continue
		}
		usedApis[cmd.API] = true
	}

	apiList := make([]string, 0, len(usedApis))
	for api := range usedApis {
		apiList = append(apiList, api)
	}
	// Sorted so generated output is byte-identical between runs. Without this,
	// Go's randomized map iteration reordered the emitted mock types every time,
	// which made a codegen drift check in CI impossible.
	sort.Strings(apiList)

	templ, err := template.New("mocks").Funcs(
		template.FuncMap{
			"Join":  strings.Join,
			"Title": capitalize,
		},
	).Parse(atMocksTemplate)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, apiList, awsatDir, "gen_mocks.go")
}

func generateAcceptanceFactory() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, specDir, func(os.FileInfo) bool { return true }, 0)
	if err != nil {
		panic(err)
	}

	finder := &findStructs{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Walk(finder, f)
		}
	}

	templ, err := template.New("acceptanceFactory").Funcs(
		template.FuncMap{
			"Title": capitalize,
		},
	).Parse(atMocksCmdBuilders)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, finder.result, awsatDir, "gen_factory.go")
}

const atMocksCmdBuilders = `/* Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// DO NOT EDIT
// This file was automatically generated with go generate
package awsat

import (
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/bootswithdefer/awless/cloud"
  awsspec "github.com/bootswithdefer/awless/aws/spec"
  "github.com/bootswithdefer/awless/logger"
)

// AcceptanceFactory builds commands whose AWS clients are routed through a Mock.
//
// SDK v2 clients are concrete structs, so they cannot be replaced with an
// interface. Instead every command is constructed from Mock.Config(), whose
// APIOptions carry a middleware that intercepts the call before it is signed or
// sent. The generated constructors already do service.NewFromConfig(cfg), so no
// SetApi call is needed.
type AcceptanceFactory struct {
	Mock   *Mock
	Logger *logger.Logger
	Graph  cloud.GraphAPI
}

func NewAcceptanceFactory(mock *Mock, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	lg := logger.DiscardLogger
	if len(l) > 0 {
		lg = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph: g, Logger: lg}
}

func (f *AcceptanceFactory) config() aws.Config {
	if f.Mock == nil {
		return aws.Config{}
	}
	return f.Mock.Config()
}

func (f *AcceptanceFactory) Build(key string) func() any {
	switch key {
		{{- range $cmdName, $cmd := . }}
		case "{{ $cmd.Action }}{{ $cmd.Entity }}":
			return func() any {
				return awsspec.New{{ $cmdName }}(f.config(), f.Graph, f.Logger)
			}
		{{- end}}
	}
	return nil
}
`

const atMocksTemplate = `/* Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// DO NOT EDIT
// This file was automatically generated with go generate
package awsat

// TODO: Acceptance mocks need reworking for AWS SDK v2.
// SDK v2 does not have iface packages, so mocks must be manually defined
// or use a different mocking strategy (e.g., httptest or interface wrappers).

{{ range $, $api := . }}

type {{ $api }}Mock struct {
  basicMock
}

{{- end }}

`
