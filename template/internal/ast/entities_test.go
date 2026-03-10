package ast

import "testing"

func TestIsInvalidEntity(t *testing.T) {
	validEntities := []string{"none", "vpc", "subnet", "instance", "securitygroup", "user", "role", "policy", "bucket", "volume", "keypair"}
	for _, e := range validEntities {
		if IsInvalidEntity(e) {
			t.Errorf("expected entity %q to be valid", e)
		}
	}

	invalidEntities := []string{"", "unknown", "VPC", "foobar", "server"}
	for _, e := range invalidEntities {
		if !IsInvalidEntity(e) {
			t.Errorf("expected entity %q to be invalid", e)
		}
	}
}
