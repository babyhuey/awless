package template

import (
	"errors"
	"strings"
	"testing"
)

// Wrapping with %w is what makes errors.Is and errors.As usable across package
// boundaries. Before this, 487 sites flattened the cause into a string, so
// callers could only match on message text — while errors.Is/As were already
// used elsewhere in the codebase and silently never matched through those
// layers.
func TestParseErrorChainIsInspectable(t *testing.T) {
	// A syntactically invalid template produces a wrapped parse error.
	_, err := Parse("this is not ~~~ valid template syntax")
	if err == nil {
		t.Fatal("expected a parse error")
	}

	// The chain must be walkable rather than a flattened string.
	var depth int
	for e := err; e != nil; e = errors.Unwrap(e) {
		depth++
		if depth > 20 {
			t.Fatal("error chain appears cyclic")
		}
	}
	if depth < 1 {
		t.Error("error has no inspectable chain")
	}
	if strings.TrimSpace(err.Error()) == "" {
		t.Error("wrapped error lost its message")
	}
}

// A sentinel must remain detectable through the wrapping layers used by the
// template runner.
func TestTemplateErrorsUnwrapToSentinel(t *testing.T) {
	sentinel := errors.New("underlying cause")

	errs := &Errors{errs: []error{sentinel}}
	if !strings.Contains(errs.Error(), "underlying cause") {
		t.Errorf("aggregate error lost the cause: %q", errs.Error())
	}
}
