package commands

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// An expired AWS credential used to print the full usage block — every global flag —
// ahead of the actual error, because cobra prints usage for any error returned from
// RunE unless told otherwise. That is a regression risk whenever a `Run:` becomes a
// `RunE:`, so the split is pinned here: cobra parses flags and validates args before
// it reaches PersistentPreRunE, so usage must survive for genuine usage mistakes and
// disappear for everything after.
func TestUsageIsSilencedForRuntimeErrorsOnly(t *testing.T) {
	runtimeErr := errors.New("failed to refresh cached credentials: ExpiredToken")

	newCmd := func() *cobra.Command {
		root := &cobra.Command{Use: "awless", SilenceErrors: true}
		root.PersistentFlags().String("aws-profile", "", "profile")
		sub := &cobra.Command{
			Use:               "list",
			Args:              cobra.MaximumNArgs(1),
			PersistentPreRunE: applyHooks(func(*cobra.Command, []string) error { return nil }),
			RunE:              func(*cobra.Command, []string) error { return runtimeErr },
		}
		root.AddCommand(sub)
		return root
	}

	for _, tc := range []struct {
		name      string
		args      []string
		wantUsage bool
	}{
		// The reported bug: a runtime failure must not be buried in usage.
		{"runtime error from RunE", []string{"list"}, false},
		// Still genuinely useful, and must not regress while fixing the above.
		{"unknown flag", []string{"list", "--nope"}, true},
		{"too many args", []string{"list", "a", "b"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newCmd()
			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)
			root.SetArgs(tc.args)

			if err := root.Execute(); err == nil {
				t.Fatal("expected an error")
			}
			gotUsage := strings.Contains(out.String(), "Usage:") ||
				strings.Contains(out.String(), "USAGE:")
			if gotUsage != tc.wantUsage {
				t.Errorf("usage printed = %v, want %v\noutput:\n%s", gotUsage, tc.wantUsage, out.String())
			}
		})
	}
}

// applyHooks must still propagate the first hook error unchanged; silencing usage is
// additive and must not swallow anything.
func TestApplyHooksPropagatesFirstError(t *testing.T) {
	sentinel := errors.New("hook failed")
	var ran int
	hooks := applyHooks(
		func(*cobra.Command, []string) error { ran++; return nil },
		func(*cobra.Command, []string) error { ran++; return sentinel },
		func(*cobra.Command, []string) error { ran++; return errors.New("must not run") },
	)

	cmd := &cobra.Command{Use: "x"}
	err := hooks(cmd, nil)
	if !errors.Is(err, sentinel) {
		t.Errorf("got %v, want %v", err, sentinel)
	}
	if ran != 2 {
		t.Errorf("ran %d hooks, want 2 (must stop at the failure)", ran)
	}
	if !cmd.SilenceUsage {
		t.Error("SilenceUsage was not set on the command")
	}
}
