package awsconfig

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

// withPromptStdin points the selectors at a scripted stdin for the duration of a test.
func withPromptStdin(t *testing.T, input string) {
	t.Helper()
	prev := promptStdin
	promptStdin = io.NopCloser(strings.NewReader(input))
	t.Cleanup(func() { promptStdin = prev })
}

// runWithTimeout fails rather than hangs. Both selectors previously looped on a stdin
// that only ever returns EOF, so a regression here is an infinite loop, and a test that
// merely called them would hang the whole suite instead of failing it.
func runWithTimeout(t *testing.T, fn func() (string, error)) (string, error) {
	t.Helper()
	type result struct {
		val string
		err error
	}
	done := make(chan result, 1)
	go func() {
		v, err := fn()
		done <- result{v, err}
	}()
	select {
	case r := <-done:
		return r.val, r.err
	case <-time.After(5 * time.Second):
		t.Fatal("selector did not return within 5s — it is looping on stdin")
		return "", nil
	}
}

func TestStdinInstanceTypeSelectorAcceptsValidInput(t *testing.T) {
	withPromptStdin(t, "t2.micro\n")

	got, err := runWithTimeout(t, StdinInstanceTypeSelector)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "t2.micro" {
		t.Errorf("got %q, want t2.micro", got)
	}
}

// An invalid entry must re-prompt rather than be accepted or abort.
func TestStdinInstanceTypeSelectorRepromptsOnInvalid(t *testing.T) {
	withPromptStdin(t, "not-a-type\nalso-wrong\nm4.large\n")

	got, err := runWithTimeout(t, StdinInstanceTypeSelector)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "m4.large" {
		t.Errorf("got %q, want m4.large", got)
	}
}

// The regression that matters: stdin closed immediately. fmt.Scan returned its error
// without assigning, so `for !isValid(...)` spun at full speed forever. Now it aborts.
func TestStdinInstanceTypeSelectorAbortsOnEOF(t *testing.T) {
	withPromptStdin(t, "")

	_, err := runWithTimeout(t, StdinInstanceTypeSelector)
	if !errors.Is(err, ErrPromptAborted) {
		t.Errorf("got %v, want ErrPromptAborted", err)
	}
}

// Same shape as EOF, but with trailing input that never validates: it must give up when
// stdin runs out rather than re-prompt against a dead reader.
func TestStdinInstanceTypeSelectorAbortsAfterExhaustingInput(t *testing.T) {
	withPromptStdin(t, "bogus\nstill-bogus\n")

	_, err := runWithTimeout(t, StdinInstanceTypeSelector)
	if !errors.Is(err, ErrPromptAborted) {
		t.Errorf("got %v, want ErrPromptAborted", err)
	}
}

// StdinRegionSelector used to call os.Exit(1) here, which would have taken the test
// process down with it. It now returns an error, which is what makes this testable.
func TestStdinRegionSelectorAbortsOnEOF(t *testing.T) {
	withPromptStdin(t, "")

	_, err := runWithTimeout(t, StdinRegionSelector)
	if !errors.Is(err, ErrPromptAborted) {
		t.Errorf("got %v, want ErrPromptAborted", err)
	}
}

func TestStdinRegionSelectorAcceptsValidRegion(t *testing.T) {
	withPromptStdin(t, "us-east-1\n")

	got, err := runWithTimeout(t, StdinRegionSelector)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "us-east-1" {
		t.Errorf("got %q, want us-east-1", got)
	}
}

func TestStdinRegionSelectorRepromptsOnInvalidRegion(t *testing.T) {
	withPromptStdin(t, "not-a-region\neu-west-1\n")

	got, err := runWithTimeout(t, StdinRegionSelector)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "eu-west-1" {
		t.Errorf("got %q, want eu-west-1", got)
	}
}

// ErrPromptAborted is what callers branch on, so it must stay distinguishable from a
// genuine read failure.
func TestErrPromptAbortedIsDistinct(t *testing.T) {
	if errors.Is(ErrPromptAborted, io.EOF) {
		t.Error("ErrPromptAborted must not match io.EOF")
	}
	wrapped := errors.New("some other failure")
	if errors.Is(wrapped, ErrPromptAborted) {
		t.Error("an unrelated error must not match ErrPromptAborted")
	}
}
