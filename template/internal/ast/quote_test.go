package ast

import "testing"

func TestQuoteStringIfNeeded(t *testing.T) {
	tcases := []struct {
		in, want string
	}{
		{"simple", "simple"},
		{"with space", "'with space'"},
		{"42", "'42'"},
		{"3.14", "'3.14'"},
		{"hello-world", "hello-world"},
		{"path/to/file", "path/to/file"},
		{"has'quote", "\"has'quote\""},
		{"a.b.c", "a.b.c"},
	}
	for _, tc := range tcases {
		if got := quoteStringIfNeeded(tc.in); got != tc.want {
			t.Errorf("quoteStringIfNeeded(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestQuote(t *testing.T) {
	tcases := []struct {
		in, want string
	}{
		{"simple", "'simple'"},
		{"has'quote", "\"has'quote\""},
		{"no quotes", "'no quotes'"},
	}
	for _, tc := range tcases {
		if got := Quote(tc.in); got != tc.want {
			t.Errorf("Quote(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
