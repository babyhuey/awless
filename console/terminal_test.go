package console

import (
	"errors"
	"testing"
)

// withTermSize substitutes the terminal size lookup for the duration of a test.
func withTermSize(t *testing.T, w, h int, err error) {
	t.Helper()
	prev := termGetSize
	termGetSize = func() (int, int, error) { return w, h, err }
	t.Cleanup(func() { termGetSize = prev })
}

func TestGetTerminalWidthAndHeight(t *testing.T) {
	withTermSize(t, 120, 40, nil)

	if got := GetTerminalWidth(); got != 120 {
		t.Errorf("width = %d, want 120", got)
	}
	if got := GetTerminalHeight(); got != 40 {
		t.Errorf("height = %d, want 40", got)
	}
}

// Redirected output is the common case in scripts and CI, where term.GetSize fails.
// Both accessors report 0 so callers can substitute their own default.
func TestGetTerminalSizeReturnsZeroOnError(t *testing.T) {
	withTermSize(t, 999, 999, errors.New("inappropriate ioctl for device"))

	if got := GetTerminalWidth(); got != 0 {
		t.Errorf("width = %d, want 0 on error", got)
	}
	if got := GetTerminalHeight(); got != 0 {
		t.Errorf("height = %d, want 0 on error", got)
	}
}

// ptySize must never hand a zero to RequestPty: a 0x0 pty leaves the remote shell with
// no usable geometry, and this is the path taken whenever awless ssh runs with output
// redirected.
func TestPtySizeSubstitutesDefaults(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		w, h                  int
		err                   error
		wantWidth, wantHeight int
	}{
		{"real size passes through", 80, 24, nil, 80, 24},
		{"error falls back on both axes", 0, 0, errors.New("not a tty"), defaultPtyDimension, defaultPtyDimension},
		{"zero width only", 0, 30, nil, defaultPtyDimension, 30},
		{"zero height only", 90, 0, nil, 90, defaultPtyDimension},
	} {
		t.Run(tc.name, func(t *testing.T) {
			withTermSize(t, tc.w, tc.h, tc.err)

			w, h := ptySize()
			if w != tc.wantWidth || h != tc.wantHeight {
				t.Errorf("ptySize() = %dx%d, want %dx%d", w, h, tc.wantWidth, tc.wantHeight)
			}
			if w == 0 || h == 0 {
				t.Error("ptySize must never return a zero dimension")
			}
		})
	}
}

// The seam must default to the real lookup so production is unaffected.
func TestTermGetSizeDefaultsToReal(t *testing.T) {
	if termGetSize == nil {
		t.Fatal("termGetSize is nil")
	}
	// Under `go test` stdout is not a tty, so this errors rather than panics. The
	// assertion is that it is callable and reports failure instead of crashing.
	if _, _, err := termGetSize(); err == nil {
		t.Log("stdout appears to be a tty; size lookup succeeded")
	}
}
