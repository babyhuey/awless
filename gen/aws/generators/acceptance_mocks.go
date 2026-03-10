package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

func generateAcceptanceMocks() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, SPEC_DIR, func(os.FileInfo) bool { return true }, 0)
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

	templ, err := template.New("mocks").Funcs(
		template.FuncMap{
			"Join":  strings.Join,
			"Title": strings.Title,
		},
	).Parse(atMocksTemplate)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, apiList, AWSAT_DIR, "gen_mocks.go")
}

func generateAcceptanceFactory() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, SPEC_DIR, func(os.FileInfo) bool { return true }, 0)
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
			"Title": strings.Title,
		},
	).Parse(atMocksCmdBuilders)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, finder.result, AWSAT_DIR, "gen_factory.go")
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
  "github.com/wallix/awless/cloud"
  awsspec "github.com/wallix/awless/aws/spec"
  "github.com/wallix/awless/logger"
)

type AcceptanceFactory struct {
	Mock   interface{}
	Logger *logger.Logger
	Graph cloud.GraphAPI
}

func NewAcceptanceFactory(mock interface{}, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	lg := logger.DiscardLogger
	if len(l) > 0 {
		lg = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph:g, Logger: lg}
}

func (f *AcceptanceFactory) Build(key string) func() interface{} {
	switch key {
		{{- range $cmdName, $cmd := . }}
		case "{{ $cmd.Action }}{{ $cmd.Entity }}":
			return func() interface{} {
				cmd := awsspec.New{{ $cmdName }}(aws.Config{}, f.Graph, f.Logger)
				// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
				_ = cmd
				return cmd
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
