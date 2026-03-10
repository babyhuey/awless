package logger

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestInteractiveInfof(t *testing.T) {
	var buf bytes.Buffer
	l := New("", 0, &buf)
	l.InteractiveInfof("loading %d%%", 50)
	out := buf.String()
	if !strings.Contains(out, "loading 50%") {
		t.Errorf("expected 'loading 50%%' in output, got: %q", out)
	}
	// Should contain carriage return for interactive clearing
	if !strings.Contains(out, "\r") {
		t.Errorf("expected carriage return in interactive output, got: %q", out)
	}
}

func TestPrepend(t *testing.T) {
	result := prepend("prefix", "a", "b")
	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
	if result[0] != "prefix" {
		t.Fatalf("expected first element 'prefix', got %v", result[0])
	}
	if result[1] != "a" {
		t.Fatalf("expected second element 'a', got %v", result[1])
	}
	if result[2] != "b" {
		t.Fatalf("expected third element 'b', got %v", result[2])
	}
}

func TestPrependSingleArg(t *testing.T) {
	result := prepend("only")
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
}

func TestFormatMultiLineErrMsgSingleLine(t *testing.T) {
	result := formatMultiLineErrMsg("simple error")
	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}
	if result[0] != "          simple error" {
		t.Fatalf("got %q", result[0])
	}
}

func TestFormatMultiLineErrMsgEmpty(t *testing.T) {
	result := formatMultiLineErrMsg("")
	if len(result) != 1 {
		t.Fatalf("expected 1 line for empty string, got %d", len(result))
	}
	if result[0] != "          " {
		t.Fatalf("expected 10 spaces, got %q", result[0])
	}
}

func TestFormatMultiLineErrMsgWithTabsAndNewlines(t *testing.T) {
	result := formatMultiLineErrMsg("line1\t\n\tline2\t\nline3")
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	for _, line := range result {
		if strings.Contains(line, "\t") {
			t.Fatalf("tabs should be removed, got: %q", line)
		}
		if !strings.HasPrefix(line, "          ") {
			t.Fatalf("expected 10-space indent, got: %q", line)
		}
	}
}

// Test the package-level functions that delegate to DefaultLogger
func TestPackageLevelInfo(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Info("test info")
	if !strings.Contains(buf.String(), "test info") {
		t.Errorf("expected 'test info', got: %q", buf.String())
	}
}

func TestPackageLevelInfof(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Infof("count=%d", 5)
	if !strings.Contains(buf.String(), "count=5") {
		t.Errorf("expected 'count=5', got: %q", buf.String())
	}
}

func TestPackageLevelError(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Error("test error")
	if !strings.Contains(buf.String(), "test error") {
		t.Errorf("expected 'test error', got: %q", buf.String())
	}
}

func TestPackageLevelErrorf(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Errorf("err %d", 42)
	if !strings.Contains(buf.String(), "err 42") {
		t.Errorf("expected 'err 42', got: %q", buf.String())
	}
}

func TestPackageLevelWarning(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Warning("test warning")
	if !strings.Contains(buf.String(), "test warning") {
		t.Errorf("expected 'test warning', got: %q", buf.String())
	}
}

func TestPackageLevelWarningf(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Warningf("warn %s", "msg")
	if !strings.Contains(buf.String(), "warn msg") {
		t.Errorf("expected 'warn msg', got: %q", buf.String())
	}
}

func TestPackageLevelVerbose(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	DefaultLogger.SetVerbose(1)
	defer func() { DefaultLogger = old }()

	Verbose("verbose msg")
	if !strings.Contains(buf.String(), "verbose msg") {
		t.Errorf("expected 'verbose msg', got: %q", buf.String())
	}
}

func TestPackageLevelVerbosef(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	DefaultLogger.SetVerbose(1)
	defer func() { DefaultLogger = old }()

	Verbosef("verbose %d", 1)
	if !strings.Contains(buf.String(), "verbose 1") {
		t.Errorf("expected 'verbose 1', got: %q", buf.String())
	}
}

func TestPackageLevelVerboseSilentByDefault(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	Verbose("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %q", buf.String())
	}
}

func TestPackageLevelExtraVerbose(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	DefaultLogger.SetVerbose(2)
	defer func() { DefaultLogger = old }()

	ExtraVerbose("extra msg")
	if !strings.Contains(buf.String(), "extra msg") {
		t.Errorf("expected 'extra msg', got: %q", buf.String())
	}
}

func TestPackageLevelExtraVerbosef(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	DefaultLogger.SetVerbose(2)
	defer func() { DefaultLogger = old }()

	ExtraVerbosef("extra %s", "formatted")
	if !strings.Contains(buf.String(), "extra formatted") {
		t.Errorf("expected 'extra formatted', got: %q", buf.String())
	}
}

func TestPackageLevelExtraVerboseSilentByDefault(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	ExtraVerbose("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %q", buf.String())
	}
}

func TestPackageLevelMultiLineError(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	MultiLineError(errors.New("line1\nline2"))
	out := buf.String()
	if !strings.Contains(out, "line1") {
		t.Errorf("expected 'line1', got: %q", out)
	}
	if !strings.Contains(out, "line2") {
		t.Errorf("expected 'line2', got: %q", out)
	}
}

func TestPackageLevelMultiLineErrorNil(t *testing.T) {
	var buf bytes.Buffer
	old := DefaultLogger
	DefaultLogger = New("", 0, &buf)
	defer func() { DefaultLogger = old }()

	MultiLineError(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for nil, got: %q", buf.String())
	}
}

func TestVerbosityLevels(t *testing.T) {
	l, buf := newTestLogger("", 0)

	// Level 0: nothing verbose
	l.SetVerbose(0)
	l.Verbose("v0")
	l.ExtraVerbose("ev0")
	if buf.Len() != 0 {
		t.Fatalf("level 0 should be silent for verbose/extra, got: %q", buf.String())
	}

	// Level 1: verbose only
	buf.Reset()
	l.SetVerbose(1)
	l.Verbose("v1")
	if !strings.Contains(buf.String(), "v1") {
		t.Fatalf("level 1 should show verbose, got: %q", buf.String())
	}
	buf.Reset()
	l.ExtraVerbose("ev1")
	if buf.Len() != 0 {
		t.Fatalf("level 1 should not show extra verbose, got: %q", buf.String())
	}

	// Level 2: both
	buf.Reset()
	l.SetVerbose(2)
	l.Verbose("v2")
	if !strings.Contains(buf.String(), "v2") {
		t.Fatalf("level 2 should show verbose, got: %q", buf.String())
	}
	buf.Reset()
	l.ExtraVerbose("ev2")
	if !strings.Contains(buf.String(), "ev2") {
		t.Fatalf("level 2 should show extra verbose, got: %q", buf.String())
	}

	// Level 3: both as well
	buf.Reset()
	l.SetVerbose(3)
	l.ExtraVerbose("ev3")
	if !strings.Contains(buf.String(), "ev3") {
		t.Fatalf("level 3 should show extra verbose, got: %q", buf.String())
	}
}
