/*
Copyright 2017 WALLIX

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

package main

import (
	"strings"
	"text/template"

	"github.com/bootswithdefer/awless/gen/aws"
)

func generateTestMocks() {
	// Build a set of APIs that need types imports based on AWSType references
	apisNeedingTypes := make(map[string]bool)
	for _, mock := range aws.Mocks() {
		for _, f := range mock.Funcs {
			if f.AWSType != "" && strings.Contains(f.AWSType, "types.") {
				apisNeedingTypes[mock.API] = true
			}
		}
	}

	templ, err := template.New("mocks").Funcs(template.FuncMap{
		"Title":            capitalize,
		"ToUpper":          strings.ToUpper,
		"Join":             strings.Join,
		"SdkModulePath":    aws.SdkModulePath,
		"NeedsTypesImport": func(api string) bool { return apisNeedingTypes[api] },
	}).Parse(mocksTempl)

	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, aws.Mocks(), servicesDir, "gen_mocks_test.go")
}

const mocksTempl = `// Auto generated implementation for the AWS cloud service

/*
Copyright 2017 WALLIX

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

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "context"

  {{- range $index, $mock := . }}
  {{ $mock.API }} "github.com/aws/aws-sdk-go-v2/service/{{ SdkModulePath $mock.API }}"
  {{- if NeedsTypesImport $mock.API }}
  {{ $mock.API }}types "github.com/aws/aws-sdk-go-v2/service/{{ SdkModulePath $mock.API }}/types"
  {{- end }}
  {{- end }}
  "github.com/bootswithdefer/awless/cloud"
)

{{ range $, $mock := . }}

type {{ $mock.Name }} struct {
	{{- range $, $func := $mock.Funcs }}
	{{- if eq $func.MockFieldType "mapslice" }}
		{{ $func.MockField}} map[string][]{{ $func.AWSType }}
	{{- else if eq $func.MockFieldType "map" }}
			{{ $func.MockField}} map[string]{{ $func.AWSType }}
	{{- else }}
		{{ $func.MockField}} []{{ $func.AWSType }}
	{{- end }}
	{{- end }}
}

func (m * {{ $mock.Name }}) Name() string {
	return ""
}

func (m * {{ $mock.Name }}) Region() string {
	return ""
}

func (m * {{ $mock.Name }}) Profile() string {
	return ""
}

func (m * {{ $mock.Name }}) Provider() string {
	return ""
}

func (m * {{ $mock.Name }}) ProviderAPI() string {
	return ""
}

func (m * {{ $mock.Name }}) ResourceTypes() []string {
	return []string{}
}

func (m * {{ $mock.Name }}) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m * {{ $mock.Name }}) IsSyncDisabled() bool {
	return false
}

func (m * {{ $mock.Name }}) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

{{ range $, $func := $mock.Funcs }}
	{{- if not $func.Manual }}
		{{- if eq $func.FuncType "list" }}
			func (m * {{ $mock.Name }}) {{ $func.APIMethod }}(ctx context.Context, input *{{ $func.Input }}, optFns ...func(*{{ $mock.API }}.Options)) (*{{ $func.Output}}, error) {
				return &{{ $func.Output}}{ {{ $func.OutputsExtractor }}: m.{{ $func.MockField}} }, nil
			}
		{{- end }}
	{{- end }}

{{ end }}

{{- end }}
`
