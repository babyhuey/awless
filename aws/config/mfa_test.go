package awsconfig

import (
	"errors"
	"testing"
)

func TestStdinMFATokenProviderAcceptsSixDigits(t *testing.T) {
	withPromptStdin(t, "123456\n")

	got, err := runWithTimeout(t, StdinMFATokenProvider)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "123456" {
		t.Errorf("got %q, want 123456", got)
	}
}

// A mistyped code should be caught here rather than becoming an AccessDenied from STS.
func TestStdinMFATokenProviderRepromptsOnMalformedInput(t *testing.T) {
	for _, bad := range []string{"12345", "1234567", "abcdef", "12 34 56", ""} {
		withPromptStdin(t, bad+"\n654321\n")

		got, err := runWithTimeout(t, StdinMFATokenProvider)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %s", bad, err)
		}
		if got != "654321" {
			t.Errorf("input %q: got %q, want the token from the second prompt", bad, got)
		}
	}
}

func TestStdinMFATokenProviderTrimsWhitespace(t *testing.T) {
	withPromptStdin(t, "  123456  \n")

	got, err := runWithTimeout(t, StdinMFATokenProvider)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got != "123456" {
		t.Errorf("got %q, want 123456", got)
	}
}

// The important case: a non-interactive stdin returns EOF immediately, so a re-prompt
// loop that does not treat EOF as an abort spins forever. This is what made two other
// prompts busy-loop before they were fixed.
func TestStdinMFATokenProviderAbortsOnEOF(t *testing.T) {
	withPromptStdin(t, "")

	_, err := runWithTimeout(t, StdinMFATokenProvider)
	if !errors.Is(err, ErrPromptAborted) {
		t.Errorf("got %v, want ErrPromptAborted", err)
	}
}

// Malformed input followed by EOF must also abort rather than loop.
func TestStdinMFATokenProviderAbortsOnEOFAfterBadInput(t *testing.T) {
	withPromptStdin(t, "nope\n")

	_, err := runWithTimeout(t, StdinMFATokenProvider)
	if !errors.Is(err, ErrPromptAborted) {
		t.Errorf("got %v, want ErrPromptAborted", err)
	}
}
