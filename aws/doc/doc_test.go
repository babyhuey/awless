package awsdoc

import (
	"strings"
	"testing"
)

func TestTemplateParamsDoc(t *testing.T) {
	tests := []struct {
		action, entity, param string
		expectedDoc           string
		expectedOk            bool
	}{
		// Known entries that exist in paramsDoc
		{
			action:      "attach",
			entity:      "alarm",
			param:       "name",
			expectedDoc: "The Name of the Alarm to update",
			expectedOk:  true,
		},
		{
			action:      "attach",
			entity:      "elasticip",
			param:       "id",
			expectedDoc: "The allocation ID",
			expectedOk:  true,
		},
		{
			action:      "attach",
			entity:      "internetgateway",
			param:       "vpc",
			expectedDoc: "The ID of the VPC",
			expectedOk:  true,
		},
		// Known action.entity exists but param does not
		{
			action:      "attach",
			entity:      "alarm",
			param:       "nonexistent-param",
			expectedDoc: "",
			expectedOk:  false,
		},
		// Action.entity combination does not exist at all
		{
			action:      "nonexistent",
			entity:      "nonexistent",
			param:       "id",
			expectedDoc: "",
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.action+"."+tt.entity+"."+tt.param, func(t *testing.T) {
			doc, ok := TemplateParamsDoc(tt.action, tt.entity, tt.param)
			if ok != tt.expectedOk {
				t.Errorf("TemplateParamsDoc(%q, %q, %q) ok = %v, want %v", tt.action, tt.entity, tt.param, ok, tt.expectedOk)
			}
			if doc != tt.expectedDoc {
				t.Errorf("TemplateParamsDoc(%q, %q, %q) doc = %q, want %q", tt.action, tt.entity, tt.param, doc, tt.expectedDoc)
			}
		})
	}
}

func TestTemplateParamsDocWithEnums(t *testing.T) {
	tests := []struct {
		name               string
		action, entity     string
		param              string
		expectContainsDoc  string
		expectContainsEnum string
		expectedOk         bool
	}{
		{
			name:               "param with both doc and enum values",
			action:             "attach",
			entity:             "policy",
			param:              "access",
			expectContainsDoc:  "Type of access to retrieve an AWS policy",
			expectContainsEnum: "readonly | full",
			expectedOk:         true,
		},
		{
			name:              "param with doc but no enum",
			action:            "attach",
			entity:            "alarm",
			param:             "name",
			expectContainsDoc: "The Name of the Alarm to update",
			expectedOk:        true,
		},
		{
			name:               "param with both doc and enum - create.instance.type",
			action:             "create",
			entity:             "instance",
			param:              "type",
			expectContainsDoc:  "The instance type",
			expectContainsEnum: "t2.nano",
			expectedOk:         true,
		},
		{
			name:              "param with neither doc nor enum",
			action:            "nonexistent",
			entity:            "nonexistent",
			param:             "id",
			expectContainsDoc: "",
			expectedOk:        false,
		},
		{
			name:               "create.database.engine has both doc and enum",
			action:             "create",
			entity:             "database",
			param:              "engine",
			expectContainsDoc:  "The name of the database engine",
			expectContainsEnum: "postgres",
			expectedOk:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, ok := TemplateParamsDocWithEnums(tt.action, tt.entity, tt.param)
			if ok != tt.expectedOk {
				t.Errorf("ok = %v, want %v", ok, tt.expectedOk)
			}
			if tt.expectContainsDoc != "" && !strings.Contains(doc, tt.expectContainsDoc) {
				t.Errorf("doc = %q, want it to contain %q", doc, tt.expectContainsDoc)
			}
			if tt.expectContainsEnum != "" && !strings.Contains(doc, tt.expectContainsEnum) {
				t.Errorf("doc = %q, want it to contain enum %q", doc, tt.expectContainsEnum)
			}
		})
	}
}

func TestAwlessCommandDefinitionsDoc(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		entity      string
		fallbackDef string
		expected    string
	}{
		{
			name:        "known command definition - copy.image",
			action:      "copy",
			entity:      "image",
			fallbackDef: "fallback",
			expected:    "Copy an EC2 image from given source region to current awless region",
		},
		{
			name:        "known command definition - create.classicloadbalancer",
			action:      "create",
			entity:      "classicloadbalancer",
			fallbackDef: "fallback",
			expected:    CommandDefinitionsDoc["create.classicloadbalancer"],
		},
		{
			name:        "unknown command returns fallback",
			action:      "nonexistent",
			entity:      "nonexistent",
			fallbackDef: "this is the fallback",
			expected:    "this is the fallback",
		},
		{
			name:        "unknown command with empty fallback",
			action:      "nonexistent",
			entity:      "nonexistent",
			fallbackDef: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AwlessCommandDefinitionsDoc(tt.action, tt.entity, tt.fallbackDef)
			if result != tt.expected {
				t.Errorf("AwlessCommandDefinitionsDoc(%q, %q, %q) = %q, want %q",
					tt.action, tt.entity, tt.fallbackDef, result, tt.expected)
			}
		})
	}
}

func TestAwlessExamplesDoc(t *testing.T) {
	tests := []struct {
		name           string
		action, entity string
		expectContains []string
		expectEmpty    bool
	}{
		{
			name:   "attach.elasticip has one example",
			action: "attach",
			entity: "elasticip",
			expectContains: []string{
				"awless attach elasticip id=eipalloc-1c517b26 instance=@redis",
			},
		},
		{
			name:   "attach.policy has multiple examples",
			action: "attach",
			entity: "policy",
			expectContains: []string{
				"awless attach policy role=MyNewRole service=ec2 access=readonly",
				"awless attach policy user=jsmith service=s3 access=readonly",
			},
		},
		{
			name:   "create.securitygroup has examples",
			action: "create",
			entity: "securitygroup",
			expectContains: []string{
				"awless create securitygroup vpc=@myvpc name=ssh-only description=ssh-access",
			},
		},
		{
			name:        "action.entity with empty examples list",
			action:      "attach",
			entity:      "instance",
			expectEmpty: true,
		},
		{
			name:        "nonexistent action.entity returns empty",
			action:      "nonexistent",
			entity:      "nonexistent",
			expectEmpty: true,
		},
		{
			name:   "delete.user has example",
			action: "delete",
			entity: "user",
			expectContains: []string{
				"awless delete user name=john",
			},
		},
		{
			name:   "update.securitygroup has multiple examples",
			action: "update",
			entity: "securitygroup",
			expectContains: []string{
				"awless update securitygroup id=@ssh-only inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=26257",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AwlessExamplesDoc(tt.action, tt.entity)
			if tt.expectEmpty {
				if result != "" {
					t.Errorf("expected empty string, got %q", result)
				}
				return
			}
			for _, expected := range tt.expectContains {
				if !strings.Contains(result, expected) {
					t.Errorf("result = %q, want it to contain %q", result, expected)
				}
			}
		})
	}
}
