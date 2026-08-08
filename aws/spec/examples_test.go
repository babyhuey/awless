package awsspec

import (
	"strings"
	"testing"

	awsdoc "github.com/bootswithdefer/awless/aws/doc"
	"github.com/bootswithdefer/awless/template/params"
)

// Every documented example must satisfy the command's own ParamsSpec.
//
// Examples appear in `-h` output for users to copy, so one the validator rejects
// is worse than no example. Several commands declare mutually exclusive required
// params — `delete instance` takes id or ids, `detach policy` takes arn or
// access+service — so an example naming every required param would be refused.
//
// Omitting a required param is fine: compilation turns it into a hole and prompts
// for it, so this applies the same normalization before validating.
func TestExamplesSatisfyParamsSpec(t *testing.T) {
	for name, def := range AWSTemplatesDefinitions {
		doc := awsdoc.AwlessExamplesDoc(def.Action, def.Entity)
		if doc == "" {
			continue // covered by TestDocForEachCommand
		}

		for _, example := range strings.Split(doc, "\n") {
			example = strings.TrimSpace(example)
			if example == "" {
				continue
			}
			keys, ok := paramKeysFromExample(example, def.Action, def.Entity)
			if !ok {
				continue // documents an alias or alternate spelling
			}
			// Mirror what compilation does: any required param the example omits
			// becomes a hole and is prompted for, so an incomplete example is
			// legitimate. What must not happen is an unknown param or a violation
			// of mutual exclusivity, which is what params.Run catches once the
			// holes are filled in.
			keys = append(keys, def.Params.Missing(keys)...)
			if err := params.Run(def.Params, keys); err != nil {
				t.Errorf("%s: example does not satisfy its own ParamsSpec: %s\n  example: %s\n  pattern: %s",
					name, err, example, def.Params)
			}
		}
	}
}

// paramKeysFromExample lists the param names an example passes.
func paramKeysFromExample(example, action, entity string) ([]string, bool) {
	rest, ok := strings.CutPrefix(example, "awless "+action+" "+entity)
	if !ok {
		return nil, false
	}

	var keys []string
	for _, field := range splitFields(strings.TrimSpace(rest)) {
		if k, _, found := strings.Cut(field, "="); found {
			keys = append(keys, k)
		}
	}
	return keys, true
}

// splitFields splits on spaces but keeps quoted values together, so a value such
// as "db creds" stays one field.
func splitFields(s string) []string {
	var fields []string
	var current strings.Builder
	var quote rune

	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote = r
		case r == ' ':
			if current.Len() > 0 {
				fields = append(fields, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		fields = append(fields, current.String())
	}
	return fields
}
