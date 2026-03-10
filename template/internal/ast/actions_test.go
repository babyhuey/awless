package ast

import "testing"

func TestIsInvalidAction(t *testing.T) {
	validActions := []string{"create", "delete", "update", "check", "start", "restart", "stop", "attach", "detach", "copy", "import", "authenticate", "none"}
	for _, a := range validActions {
		if IsInvalidAction(a) {
			t.Errorf("expected action %q to be valid", a)
		}
	}

	invalidActions := []string{"", "unknown", "fly", "CREATE", "Deploy"}
	for _, a := range invalidActions {
		if !IsInvalidAction(a) {
			t.Errorf("expected action %q to be invalid", a)
		}
	}
}
