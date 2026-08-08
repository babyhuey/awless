package web

import (
	"html/template"
	"strings"
	"testing"
)

// The resource page is an html/template over web.Resource, so renaming a field
// without updating the template fails only when the page is served. Nothing else
// in the suite executes this template, which is how the Id -> ID rename slipped
// through.
func TestShowResourceTemplateMatchesResourceFields(t *testing.T) {
	tpl, err := template.New("show").Parse(showResourceTpl)
	if err != nil {
		t.Fatalf("template does not parse: %s", err)
	}

	res := &Resource{
		ID:         "i-1234",
		Type:       "instance",
		Properties: map[string]any{"Name": "redis"},
		Parents:    []*Resource{{ID: "vpc-1", Type: "vpc"}},
		Children:   []*Resource{{ID: "sub-1", Type: "subnet"}},
		DependsOn:  []*Resource{{ID: "sg-1", Type: "securitygroup"}},
		AppliesOn:  []*Resource{{ID: "vol-1", Type: "volume"}},
	}

	var out strings.Builder
	if err := tpl.Execute(&out, res); err != nil {
		t.Fatalf("executing the resource template: %s", err)
	}

	// Every nested list is rendered, so a field missed in any branch shows up here.
	for _, want := range []string{"i-1234", "vpc-1", "sub-1", "sg-1", "vol-1"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("expected %q in the rendered page", want)
		}
	}
}
