package main

import (
	"sort"
	"text/template"
)

// The template parser validates entities against a map that was maintained by
// hand, independently of the commands themselves. A command could therefore be
// fully registered — visible in `awless <action> <entity> -h`, with working param
// validation — while every template or one-liner using it failed at parse time
// with `unknown entity`. That is how Secrets Manager and SSM shipped initially.
//
// Generating the map from the same `entity:` struct tags the rest of the
// generators already read removes the possibility of drift, and the codegen CI
// job now fails if the committed file falls out of date.
func generateTemplateEntities() {
	cmds := loadCommandStructs()

	seen := make(map[string]bool)
	for _, cmd := range cmds {
		seen[cmd.Entity] = true
	}

	// Entities kept parseable despite having no command.
	//
	// `awless log` re-parses persisted template lines, so an entity that is
	// dropped here breaks history that already exists in users' ~/.awless. These
	// must stay until that is no longer a concern:
	//
	//   none              the parser's placeholder for a statement with no entity
	//   container         referenced by template/revert.go; a cloud.Container
	//                     resource exists and is list-only
	//   containerservice  cloud.ContainerService, list-only
	for _, extra := range []string{"none", "container", "containerservice"} {
		seen[extra] = true
	}

	entities := make([]string, 0, len(seen))
	for e := range seen {
		entities = append(entities, e)
	}
	sort.Strings(entities)

	templ := template.Must(template.New("entities").Parse(entitiesTemplate))
	writeTemplateToFile(templ, entities, templateASTDir, "gen_entities.go")
}

const entitiesTemplate = `/* Copyright 2017 WALLIX

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
package ast

// entities is the set the template parser accepts. Derived from the ` + "`entity:`" + `
// struct tags in aws/spec, plus a small set of entities kept parseable so that
// template history already written to ~/.awless still loads. See
// gen/aws/generators/entities.go.
var entities = map[Entity]struct{}{
{{- range . }}
	"{{ . }}": {},
{{- end }}
}
`
