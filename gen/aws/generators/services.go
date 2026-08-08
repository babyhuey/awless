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

func generateServicesFuncs() {
	// Build a set of APIs that need types imports based on AWSType references
	apisNeedingTypes := make(map[string]bool)
	for _, def := range aws.FetchersDefs {
		for _, f := range def.Fetchers {
			if idx := strings.Index(f.AWSType, "types."); idx > 0 {
				apisNeedingTypes[f.AWSType[:idx]] = true
			}
		}
	}

	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":            capitalize,
		"ToUpper":          strings.ToUpper,
		"Join":             strings.Join,
		"SdkModulePath":    aws.SdkModulePath,
		"NeedsTypesImport": func(api string) bool { return apisNeedingTypes[api] },
	}).Parse(servicesTempl)

	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, aws.FetchersDefs, servicesDir, "gen_services.go")
}

const servicesTempl = `// Auto generated implementation for the AWS cloud service

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
  "errors"
	"sync"

  "github.com/aws/aws-sdk-go-v2/aws"
  {{- range $index, $service := . }}
  {{- range $, $api := $service.API }}
  {{ $api }} "github.com/aws/aws-sdk-go-v2/service/{{ SdkModulePath $api }}"
  {{- if NeedsTypesImport $api }}
  {{ $api }}types "github.com/aws/aws-sdk-go-v2/service/{{ SdkModulePath $api }}/types"
  {{- end }}
  {{- end }}
  {{- end }}
  "github.com/aws/smithy-go"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/fetch"
	awsfetch "github.com/bootswithdefer/awless/aws/fetch"
	tstore "github.com/bootswithdefer/triplestore"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	{{- range $index, $service := . }}
  "{{ $service.Name }}",
  {{- end }}
}

var ResourceTypes = []string {
{{- range $index, $service := . }}
    {{- range $idx, $fetcher := $service.Fetchers }}
      "{{ $fetcher.ResourceType }}",
    {{- end }}
{{- end }}
}

var ServicePerAPI = map[string]string {
{{- range $index, $service := . }}
{{- range $, $api := $service.API }}
  "{{ $api }}": "{{ $service.Name }}",
{{- end }}
{{- end }}
}

var ServicePerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $service.Name }}",
  {{- end }}
{{- end }}
}

var APIPerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $fetcher.API }}",
  {{- end }}
{{- end }}
}

{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
	fetcher fetch.Fetcher
  region, profile string
	config map[string]any
	log *logger.Logger
	{{- range $, $api := $service.API }}
		{{ Title $api }}Client *{{ $api }}.Client
	{{- end }}
}

func New{{ Title $service.Name }}(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
  {{- if $service.Global }}
	region := "global"
	{{- else}}
	region := cfg.Region
	{{- end}}

	{{- range $, $api := $service.API }}
		{{ $api }}Client := {{ $api }}.NewFromConfig(cfg)
	{{- end }}

	fetchConfig := awsfetch.NewConfig(
		{{- range $, $api := $service.API }}
			{{ $api }}Client,
		{{- end }}
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &{{ Title $service.Name }}{
	{{- range $, $api := $service.API }}
		{{ Title $api }}Client: {{ $api }}Client,
	{{- end }}
		fetcher: fetch.NewFetcher(awsfetch.Build{{ Title $service.Name }}FetchFuncs(fetchConfig)),
		config: extraConf,
		region: region,
		profile: profile,
		log: log,
  }
}

func (s *{{ Title $service.Name }}) Name() string {
  return "{{ $service.Name }}"
}

func (s *{{ Title $service.Name }}) Region() string {
  return s.region
}

func (s *{{ Title $service.Name }}) Profile() string {
  return s.profile
}

func (s *{{ Title $service.Name }}) ResourceTypes() []string {
	return []string{
	{{- range $index, $fetcher := $service.Fetchers }}
		"{{ $fetcher.ResourceType }}",
	{{- end }}
	}
}

func (s *{{ Title $service.Name }}) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

  gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()
	
	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup

	{{- range $index, $fetcher := $service.Fetchers }}
	if getBool(s.config, "aws.{{ $service.Name }}.{{ $fetcher.ResourceType }}.sync", true) {
		list, err := s.fetcher.Get("{{ $fetcher.ResourceType }}_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]{{ $fetcher.AWSType }}); !ok {
			return gph, errors.New("cannot cast to '[]{{ $fetcher.AWSType }}' type from fetch context")
		}
		for _, r := range list.([]{{ $fetcher.AWSType }}) {
			for _, fn := range addParentsFns["{{ $fetcher.ResourceType }}"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *{{ $fetcher.AWSType }}) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
  {{- end }}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
				allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *{{ Title $service.Name }}) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
  return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *{{ Title $service.Name }}) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.{{ $service.Name }}.sync", true)
}

{{ end }}`
