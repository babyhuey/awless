package template

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const secret = "sup3rs3cr3t-P@ssw0rd"

// The persisted template log must never contain a plaintext secret. This
// asserts the requirement (secret absent) rather than the implementation, so it
// still holds if the redaction mechanism changes.
func TestMarshalJSONRedactsSensitiveParams(t *testing.T) {
	tcases := []struct {
		desc string
		tpl  string
	}{
		{desc: "create loginprofile", tpl: "create loginprofile username=jsmith password=" + secret},
		{desc: "update loginprofile", tpl: "update loginprofile username=jsmith password=" + secret},
		{desc: "create database", tpl: "create database id=db1 password=" + secret},
		{desc: "quoted value", tpl: `create loginprofile username=jsmith password="` + secret + `"`},
		{desc: "multiple commands", tpl: "create loginprofile username=a password=" + secret + "\ncreate vpc cidr=10.0.0.0/16"},
	}

	for _, tc := range tcases {
		t.Run(tc.desc, func(t *testing.T) {
			tpl, err := Parse(tc.tpl)
			if err != nil {
				t.Fatalf("parsing %q: %s", tc.tpl, err)
			}

			texec := &TemplateExecution{Template: tpl}
			b, err := texec.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if strings.Contains(string(b), secret) {
				t.Errorf("secret leaked into persisted JSON:\n%s", b)
			}
		})
	}
}

func TestMarshalJSONRedactsFillers(t *testing.T) {
	tpl, err := Parse("create loginprofile username=jsmith password={pwd}")
	if err != nil {
		t.Fatal(err)
	}

	texec := &TemplateExecution{
		Template: tpl,
		Fillers:  map[string]interface{}{"password": secret, "username": "jsmith"},
	}

	b, err := texec.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), secret) {
		t.Errorf("secret leaked into persisted fillers:\n%s", b)
	}

	// Non-sensitive fillers must survive untouched.
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	fillers, ok := out["fillers"].(map[string]interface{})
	if !ok {
		t.Fatalf("fillers missing or wrong type: %#v", out["fillers"])
	}
	if got := fillers["username"]; got != "jsmith" {
		t.Errorf("non-sensitive filler altered: got %v, want jsmith", got)
	}
}

// Redacted command lines are re-parsed when the log is read back, so they must
// remain valid template syntax.
func TestRedactedLineRoundTripsThroughParser(t *testing.T) {
	tpl, err := Parse("create loginprofile username=jsmith password=" + secret)
	if err != nil {
		t.Fatal(err)
	}

	texec := &TemplateExecution{Template: tpl}
	b, err := texec.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	reloaded := &TemplateExecution{}
	if err := reloaded.UnmarshalJSON(b); err != nil {
		t.Fatalf("redacted line did not survive a round trip: %s", err)
	}

	var count int
	for _, cmd := range reloaded.CommandNodesIterator() {
		count++
		// Params reload as ast.InterfaceNode wrappers, so compare rendered form.
		if got := fmt.Sprint(cmd.ParamNodes["username"]); got != "jsmith" {
			t.Errorf("non-sensitive param altered: got %q, want jsmith", got)
		}
		got := fmt.Sprint(cmd.ParamNodes["password"])
		if got == secret {
			t.Error("secret survived the round trip")
		}
		if got != "*****" {
			t.Errorf("password not redacted on reload: got %q", got)
		}
	}
	if count != 1 {
		t.Errorf("expected 1 command after reload, got %d", count)
	}
}

// Display and execution paths must be unaffected: redaction happens only at the
// persistence boundary.
func TestStringIsNotRedacted(t *testing.T) {
	tpl, err := Parse("create loginprofile username=jsmith password=" + secret)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(tpl.String(), secret) {
		t.Error("String() should not redact; only persistence and logging should")
	}
	if strings.Contains(tpl.StringRedacted(), secret) {
		t.Error("StringRedacted() must redact")
	}
}
