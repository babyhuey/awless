package inspect

import (
	"testing"
)

func TestInspectorsRegisterHasExpectedNames(t *testing.T) {
	expectedNames := []string{"pricer", "bucket_sizer", "port_scanner", "open_buckets"}

	for _, name := range expectedNames {
		if _, ok := InspectorsRegister[name]; !ok {
			t.Errorf("expected inspector '%s' to be registered, but it was not found", name)
		}
	}

	if len(InspectorsRegister) != len(expectedNames) {
		t.Errorf("expected %d registered inspectors, got %d", len(expectedNames), len(InspectorsRegister))
	}
}

func TestRegisteredInspectorsHaveNonEmptyNames(t *testing.T) {
	for key, insp := range InspectorsRegister {
		name := insp.Name()
		if name == "" {
			t.Errorf("inspector registered under key '%s' has an empty Name()", key)
		}
		if name != key {
			t.Errorf("inspector Name() '%s' does not match its registration key '%s'", name, key)
		}
	}
}
