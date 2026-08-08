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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"sort"

	"github.com/bootswithdefer/awless/gen/aws"
)

func loadCommandStructs() map[string]cmdData {
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
	return finder.result
}

func uniqueAPIs(cmds map[string]cmdData) []string {
	seen := make(map[string]bool)
	var apis []string
	for _, cmd := range cmds {
		if !seen[cmd.API] {
			seen[cmd.API] = true
			apis = append(apis, cmd.API)
		}
	}
	sort.Strings(apis)
	return apis
}

func generateCommands() {
	cmdsData := loadCommandStructs()
	templ, err := template.New("cmdRuns").Funcs(
		template.FuncMap{
			"SdkModulePath": aws.SdkModulePath,
			"UniqueAPIs":    uniqueAPIs,
		},
	).Parse(cmdRuns)
	if err != nil {
		panic(err)
	}
	writeTemplateToFile(templ, cmdsData, specDir, "gen_runs.go")

	templ, err = template.New("cmdInits").Parse(cmdInits)
	if err != nil {
		panic(err)
	}
	writeTemplateToFile(templ, cmdsData, specDir, "gen_inits.go")

	templ, err = template.New("templates_definitions").Funcs(
		template.FuncMap{
			"BuildSupportedActions": BuildSupportedActions,
		},
	).Parse(cmdsDefinitions)
	if err != nil {
		panic(err)
	}
	writeTemplateToFile(templ, cmdsData, specDir, "gen_cmds_defs.go")
}

type cmdData struct {
	Action, Entity, API, Call, Input, Output string
	Params                                   []templateParam
	RequiredParamsKey                        []string
	ExtrasParamsKey                          []string
	HasRequiredParams                        bool
	HasDryRun                                bool
	GenDryRun                                bool
}

type templateParam struct {
	Name       string
	AwsField   string
	IsRequired bool
}

type findStructs struct {
	result map[string]cmdData
}

func (v *findStructs) Visit(node ast.Node) (w ast.Visitor) {
	if v.result == nil {
		v.result = make(map[string]cmdData)
	}
	if typ, ok := node.(*ast.TypeSpec); ok {
		if s, isStruct := typ.Type.(*ast.StructType); isStruct {
			var cmd *cmdData
			var params []templateParam
			for _, f := range s.Fields.List {
				if tag := f.Tag; tag != nil && strings.Contains(tag.Value, "awsAPI") {
					extractedCmd := extractCmdData(tag.Value)
					cmd = &extractedCmd
					continue
				}
				if tag := f.Tag; tag != nil && strings.Contains(tag.Value, "templateName") {
					params = append(params, extractParam(tag.Value))
				}
			}
			if cmd != nil {
				if len(params) > 0 {
					sort.Slice(params, func(i, j int) bool {
						return params[i].Name < params[j].Name
					})
					cmd.Params = params
				}
				for _, p := range cmd.Params {
					if p.IsRequired {
						cmd.HasRequiredParams = true
						cmd.RequiredParamsKey = append(cmd.RequiredParamsKey, p.Name)
					} else {
						cmd.ExtrasParamsKey = append(cmd.ExtrasParamsKey, p.Name)
					}
				}
				v.result[typ.Name.Name] = *cmd
			}
		}
	}
	return v
}

func extractCmdData(s string) (t cmdData) {
	tags := extractTags(s)
	if v, ok := tags["action"]; ok {
		t.Action = v
	}
	if v, ok := tags["entity"]; ok {
		t.Entity = v
	}
	if v, ok := tags["awsAPI"]; ok {
		t.API = v
	}
	if v, ok := tags["awsCall"]; ok {
		t.Call = v
	}
	if v, ok := tags["awsInput"]; ok {
		t.Input = v
	}
	if v, ok := tags["awsOutput"]; ok {
		t.Output = v
	}
	if v, ok := tags["awsDryRun"]; ok {
		t.HasDryRun = true
		t.GenDryRun = true
		if strings.ToLower(v) == "manual" {
			t.GenDryRun = false
		}
	}
	return
}

func extractParam(s string) (p templateParam) {
	tags := extractTags(s)
	if v, ok := tags["templateName"]; ok {
		p.Name = v
	}
	if _, ok := tags["required"]; ok {
		p.IsRequired = true
	}
	if v, ok := tags["awsName"]; ok {
		p.AwsField = v
	}
	return
}

func BuildSupportedActions(cmds map[string]cmdData) map[string][]string {
	supportedActions := make(map[string][]string)
	for _, cmd := range cmds {
		supportedActions[cmd.Action] = append(supportedActions[cmd.Action], cmd.Entity)
	}
	for _, entities := range supportedActions {
		sort.Slice(entities, func(i, j int) bool {
			return entities[i] < entities[j]
		})
	}
	return supportedActions
}

func extractTags(s string) map[string]string {
	splits := strings.Split(s[1:len(s)-1], " ")
	tags := make(map[string]string)
	for _, e := range splits {
		el := strings.Split(e, ":")
		if len(el) > 1 {
			if len(el[1]) < 2 || el[1][0] != '"' || el[1][len(el[1])-1] != '"' {
				panic(fmt.Sprintf("malformed tag: '%s':'%s'", el[0], el[1]))
			}
			tags[el[0]] = el[1][1 : len(el[1])-1]
		}
	}
	return tags
}

const cmdRuns = `/* Copyright 2017 WALLIX

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
package awsspec

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	{{- range $, $api := UniqueAPIs . }}
	{{ $api }} "github.com/aws/aws-sdk-go-v2/service/{{ SdkModulePath $api }}"
	{{- end }}
	"github.com/aws/smithy-go"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
)

{{ range $cmdName, $tag := . }}
func New{{ $cmdName }}(cfg aws.Config, g cloud.GraphAPI, l ...*logger.Logger) *{{ $cmdName }}{
	cmd := new({{ $cmdName }})
	if len(l) > 0 {
		cmd.logger = l[0]
	} else {
		cmd.logger = logger.DiscardLogger
	}
	if cfg.Region != "" {
		cmd.api = {{ $tag.API }}.NewFromConfig(cfg)
	}
	cmd.graph = g
	return cmd
}

func (cmd *{{ $cmdName }}) Run(renv env.Running, params map[string]any) (any, error) {
	if renv.IsDryRun() {
		return cmd.dryRun(renv, params)
	}
	return cmd.run(renv, params)
}

func (cmd *{{ $cmdName }}) run(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}
	
	if v, ok := implementsBeforeRun(cmd); ok {
		if brErr := v.BeforeRun(renv); brErr != nil {
			return nil, fmt.Errorf("before run: %s", brErr)
		}
	}
	
	{{ if $tag.Call }}
	input := &{{ $tag.Input }}{}
	if err := structInjector(cmd, input, renv) ; err != nil {
		return nil, fmt.Errorf("cannot inject in {{ $tag.Input }}: %w", err)
	}
	if v, ok := implementsInputPostProcessor(cmd); ok {
		v.PostProcessInput(input)
	}
	start := time.Now()
	output, err := cmd.api.{{ $tag.Call }}(renv.RequestContext(), input)
	renv.Log().ExtraVerbosef("{{ $tag.API }}.{{ $tag.Call }} call took %s", time.Since(start))
	if err != nil {
		return nil, decorateAWSError(err)
	}
	{{- else }}
	
	output, err := cmd.ManualRun(renv)
	if err != nil {
		return nil, decorateAWSError(err)
	}
	{{- end }}
	
	var extracted any
	if v, ok := implementsResultExtractor(cmd); ok {
		if output != nil {
			extracted = v.ExtractResult(output)
		} else {
			renv.Log().Warning("{{ $tag.Action }} {{ $tag.Entity }}: AWS command returned nil output")
		}
	}
	
	if extracted != nil {
		renv.Log().Verbosef("{{ $tag.Action }} {{ $tag.Entity }} '%s' done", extracted)
	} else {
		renv.Log().Verbose("{{ $tag.Action }} {{ $tag.Entity }} done")
	}

	if v, ok := implementsAfterRun(cmd); ok {
		if brErr := v.AfterRun(renv, output); brErr != nil {
			return nil, fmt.Errorf("after run: %s", brErr)
		}
	}

	return extracted, nil
}

{{ if $tag.HasDryRun }}
	{{ if $tag.GenDryRun }}
	func (cmd *{{ $cmdName }}) dryRun(renv env.Running, params map[string]any) (any, error) {
		if err := cmd.inject(params); err != nil {
			return nil, fmt.Errorf("cannot set params on command struct: %w", err)
		}

		input := &{{ $tag.Input }}{}
		input.DryRun = aws.Bool(true)
		if err := structInjector(cmd, input, renv) ; err != nil {
			return nil, fmt.Errorf("cannot inject in {{ $tag.Input }}: %w", err)
		}
		if v, ok := implementsInputPostProcessor(cmd); ok {
			v.PostProcessInput(input)
		}

		start := time.Now()
		_, err := cmd.api.{{ $tag.Call }}(renv.RequestContext(), input);
		var ae smithy.APIError
		if errors.As(err, &ae) {
			switch code := ae.ErrorCode(); {
			case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(ae.ErrorMessage(), "Invalid IAM Instance Profile name"):
				renv.Log().ExtraVerbosef("dry run: {{ $tag.API }}.{{ $tag.Call }} call took %s", time.Since(start))
				renv.Log().Verbose("dry run: {{ $tag.Action }} {{ $tag.Entity }} ok")
				return fakeDryRunID("{{ $tag.Entity }}"), nil
			}
		}

		return nil, err
	}
	{{- end }}
{{- else }}
func (cmd *{{ $cmdName }}) dryRun(renv env.Running, params map[string]any) (any, error) {
	return fakeDryRunID("{{ $tag.Entity }}"), nil
}
{{- end }}

func (cmd *{{ $cmdName }}) inject(params map[string]any) error {
	return structSetter(cmd, params)
}
{{ end }}
`

const cmdInits = `/* Copyright 2017 WALLIX

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
package awsspec

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
)

type Factory interface {
	Build(key string) func() any
}

var CommandFactory Factory

var MockAWSSessionFactory = &AWSFactory{
	Log: logger.DiscardLogger,
	Cfg: aws.Config{},
}

type AWSFactory struct {
	Log   *logger.Logger
	Cfg   aws.Config
	Graph cloud.GraphAPI
}

func (f *AWSFactory) Build(key string) func() any {
	switch key {
	{{- range $cmdName, $tag := . }}
	case "{{ $tag.Action }}{{ $tag.Entity }}":
		return func() any { return New{{ $cmdName }}(f.Cfg, f.Graph, f.Log) }
	{{- end}}
	}
	return nil
}

var (
	{{- range $cmdName, $tag := . }}
	_ command = &{{ $cmdName }}{}
	{{- end }}
)
`

const cmdsDefinitions = `/* Copyright 2017 WALLIX

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
package awsspec

var APIPerTemplateDefName = map[string]string {
{{- range $, $cmd := . }}
  "{{ $cmd.Action }}{{ $cmd.Entity }}": "{{ $cmd.API }}",
{{- end }}
}

var AWSTemplatesDefinitions = map[string]Definition{
{{- range $cmdName, $cmd := . }}
	"{{ $cmd.Action }}{{ $cmd.Entity }}": Definition{
			Action: "{{ $cmd.Action }}",
			Entity: "{{ $cmd.Entity }}",
			API: "{{ $cmd.API }}",
			Params: new({{ $cmdName }}).ParamsSpec().Rule(),
		},
{{- end }}
}

var DriverSupportedActions = map[string][]string{
	{{- range $action, $entities := BuildSupportedActions . }}
		"{{ $action }}" : []string{ {{- range $entity := $entities }}"{{ $entity }}", {{- end}} },
	{{- end }}
}


`
