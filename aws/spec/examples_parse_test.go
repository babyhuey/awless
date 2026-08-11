package awsspec

import (
	"strings"
	"testing"

	awsdoc "github.com/bootswithdefer/awless/aws/doc"
	"github.com/bootswithdefer/awless/template"
)

// Every documented example must parse with the real template grammar.
//
// TestExamplesSatisfyParamsSpec validates the param *names* an example uses, but it
// splits fields itself rather than parsing, so an example that no user could actually run
// still passed. That is how a rule pattern documented as
// `pattern="{\"source\":[\"aws.ec2\"]}"` got in: the grammar has no escape for a double
// quote inside a double-quoted value, so it must be single-quoted.
//
// Examples are compiled the same way a user's command line is, so this is the check that
// a copy-paste works.
func TestExamplesParse(t *testing.T) {
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
			// The examples are written as full command lines; the template language
			// is the same text without the binary name.
			line, ok := strings.CutPrefix(example, "awless ")
			if !ok {
				t.Errorf("%s: example does not start with 'awless ': %s", name, example)
				continue
			}
			if _, err := template.Parse(line); err != nil {
				t.Errorf("%s: documented example does not parse: %s\n  example: %s",
					name, err, example)
			}
		}
	}
}
