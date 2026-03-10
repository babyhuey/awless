package logger

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func newTestLogger(prefix string, flag int) (*Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	l := New(prefix, flag, &buf)
	return l, &buf
}

func TestNew(t *testing.T) {
	l, buf := newTestLogger("myprefix", 0)
	l.Info("hello")
	out := buf.String()
	if !strings.Contains(out, "myprefix") {
		t.Errorf("expected output to contain prefix 'myprefix', got: %q", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected output to contain 'hello', got: %q", out)
	}
}

func TestNewDefaultsToStderr(t *testing.T) {
	// When no writer is provided, it should not panic
	l := New("", 0)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInfoOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Info("test message")
	out := buf.String()
	if !strings.Contains(out, "info") {
		t.Errorf("expected Info output to contain 'info', got: %q", out)
	}
	if !strings.Contains(out, "test message") {
		t.Errorf("expected Info output to contain 'test message', got: %q", out)
	}
}

func TestInfofOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Infof("count=%d", 42)
	out := buf.String()
	if !strings.Contains(out, "count=42") {
		t.Errorf("expected Infof output to contain 'count=42', got: %q", out)
	}
}

func TestWarningOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Warning("warn msg")
	out := buf.String()
	if !strings.Contains(out, "warning") {
		t.Errorf("expected Warning output to contain 'warning', got: %q", out)
	}
	if !strings.Contains(out, "warn msg") {
		t.Errorf("expected Warning output to contain 'warn msg', got: %q", out)
	}
}

func TestWarningfOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Warningf("val=%s", "abc")
	out := buf.String()
	if !strings.Contains(out, "val=abc") {
		t.Errorf("expected Warningf output to contain 'val=abc', got: %q", out)
	}
}

func TestErrorOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Error("err msg")
	out := buf.String()
	if !strings.Contains(out, "error") {
		t.Errorf("expected Error output to contain 'error', got: %q", out)
	}
	if !strings.Contains(out, "err msg") {
		t.Errorf("expected Error output to contain 'err msg', got: %q", out)
	}
}

func TestErrorfOutput(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Errorf("code=%d", 500)
	out := buf.String()
	if !strings.Contains(out, "code=500") {
		t.Errorf("expected Errorf output to contain 'code=500', got: %q", out)
	}
}

func TestVerboseNotShownAtDefaultLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Verbose("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output at default verbosity, got: %q", buf.String())
	}
}

func TestVerbosefNotShownAtDefaultLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Verbosef("should not appear %d", 1)
	if buf.Len() != 0 {
		t.Errorf("expected no output at default verbosity, got: %q", buf.String())
	}
}

func TestVerboseShownAtVerboseLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(1)
	l.Verbose("verbose msg")
	out := buf.String()
	if !strings.Contains(out, "verbose msg") {
		t.Errorf("expected verbose output, got: %q", out)
	}
}

func TestVerbosefShownAtVerboseLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(1)
	l.Verbosef("verbose %s", "formatted")
	out := buf.String()
	if !strings.Contains(out, "verbose formatted") {
		t.Errorf("expected verbose formatted output, got: %q", out)
	}
}

func TestExtraVerboseNotShownAtDefaultLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.ExtraVerbose("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output at default verbosity, got: %q", buf.String())
	}
}

func TestExtraVerbosefNotShownAtDefaultLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.ExtraVerbosef("should not appear %d", 1)
	if buf.Len() != 0 {
		t.Errorf("expected no output at default verbosity, got: %q", buf.String())
	}
}

func TestExtraVerboseNotShownAtVerboseLevel1(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(1)
	l.ExtraVerbose("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output at verbosity level 1, got: %q", buf.String())
	}
}

func TestExtraVerboseShownAtVerboseLevel2(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(2)
	l.ExtraVerbose("extra msg")
	out := buf.String()
	if !strings.Contains(out, "extra msg") {
		t.Errorf("expected extra verbose output, got: %q", out)
	}
}

func TestExtraVerbosefShownAtVerboseLevel2(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(2)
	l.ExtraVerbosef("extra %s", "formatted")
	out := buf.String()
	if !strings.Contains(out, "extra formatted") {
		t.Errorf("expected extra verbose formatted output, got: %q", out)
	}
}

func TestVerboseAlsoShownAtLevel2(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.SetVerbose(2)
	l.Verbose("verbose at level 2")
	out := buf.String()
	if !strings.Contains(out, "verbose at level 2") {
		t.Errorf("expected verbose output at level 2, got: %q", out)
	}
}

func TestSetVerboseChangesLevel(t *testing.T) {
	l, buf := newTestLogger("", 0)

	// Initially no verbose output
	l.Verbose("hidden")
	if buf.Len() != 0 {
		t.Fatal("expected no output before SetVerbose")
	}

	// Enable verbose
	l.SetVerbose(1)
	l.Verbose("visible")
	if !strings.Contains(buf.String(), "visible") {
		t.Fatal("expected verbose output after SetVerbose(1)")
	}

	// Reset
	buf.Reset()
	l.SetVerbose(0)
	l.Verbose("hidden again")
	if buf.Len() != 0 {
		t.Fatal("expected no output after SetVerbose(0)")
	}
}

func TestDiscardLoggerDiscardsEverything(t *testing.T) {
	l := New("", 0, io.Discard)
	l.SetVerbose(2)

	// None of these should panic or produce output
	l.Info("discarded")
	l.Infof("discarded %d", 1)
	l.Warning("discarded")
	l.Warningf("discarded %d", 1)
	l.Error("discarded")
	l.Errorf("discarded %d", 1)
	l.Verbose("discarded")
	l.Verbosef("discarded %d", 1)
	l.ExtraVerbose("discarded")
	l.ExtraVerbosef("discarded %d", 1)
	l.MultiLineError(errors.New("discarded"))
	l.Println()
}

func TestMultiLineErrorNilError(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.MultiLineError(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for nil error, got: %q", buf.String())
	}
}

func TestMultiLineErrorSingleLine(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.MultiLineError(errors.New("single line error"))
	out := buf.String()
	if !strings.Contains(out, "single line error") {
		t.Errorf("expected output to contain error message, got: %q", out)
	}
}

func TestMultiLineErrorMultipleLines(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.MultiLineError(errors.New("line1\nline2\nline3"))
	out := buf.String()
	if !strings.Contains(out, "line1") {
		t.Errorf("expected output to contain 'line1', got: %q", out)
	}
	if !strings.Contains(out, "line2") {
		t.Errorf("expected output to contain 'line2', got: %q", out)
	}
	if !strings.Contains(out, "line3") {
		t.Errorf("expected output to contain 'line3', got: %q", out)
	}
}

func TestMultiLineErrorRemovesTabs(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.MultiLineError(errors.New("has\ttab"))
	out := buf.String()
	if strings.Contains(out, "\t") {
		t.Errorf("expected tabs to be removed, got: %q", out)
	}
	if !strings.Contains(out, "hastab") {
		t.Errorf("expected 'hastab' (tab removed), got: %q", out)
	}
}

func TestMultiLineErrorIndentation(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.MultiLineError(errors.New("indented"))
	out := buf.String()
	// formatMultiLineErrMsg adds 10 spaces of indentation
	if !strings.Contains(out, "          indented") {
		t.Errorf("expected 10-space indentation, got: %q", out)
	}
}

func TestFormatMultiLineErrMsg(t *testing.T) {
	result := formatMultiLineErrMsg("a\nb\tc")
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
	if result[0] != "          a" {
		t.Errorf("expected '          a', got %q", result[0])
	}
	// Tab should be removed
	if result[1] != "          bc" {
		t.Errorf("expected '          bc', got %q", result[1])
	}
}

func TestVerboseFAndExtraVerboseFConstants(t *testing.T) {
	if VerboseF != 1 {
		t.Errorf("expected VerboseF == 1, got %d", VerboseF)
	}
	if ExtraVerboseF != 2 {
		t.Errorf("expected ExtraVerboseF == 2, got %d", ExtraVerboseF)
	}
}

func TestPrintln(t *testing.T) {
	l, buf := newTestLogger("", 0)
	l.Println()
	out := buf.String()
	if out != "\n" {
		t.Errorf("expected single newline from Println, got: %q", out)
	}
}
