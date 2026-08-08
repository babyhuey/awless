package awsservices

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// A monitor whose requests all share a timestamp gives maxDuration == 0, which was
// divided by when scaling each bar. Found via gosec's G115 report on the uint
// conversions here.
func TestDisplayWithZeroDuration(t *testing.T) {
	now := time.Now()
	n := &NetworkMonitor{
		requests: map[string]*req{
			"a": {name: "ec2.DescribeInstances", from: now, to: now},
		},
	}

	var buf bytes.Buffer
	n.DisplayStats(&buf) // must not panic with a division by zero

	if !strings.Contains(buf.String(), "ec2.DescribeInstances") {
		t.Errorf("expected the request name in the output, got %q", buf.String())
	}
}

// A function name long enough to exceed the terminal width made maxwidth wrap to a
// huge uint, and strings.Repeat then tried to allocate it.
func TestDisplayWithNameWiderThanTerminal(t *testing.T) {
	now := time.Now()
	n := &NetworkMonitor{
		requests: map[string]*req{
			"a": {name: strings.Repeat("averyLongServiceCallName", 20), from: now, to: now.Add(time.Second)},
		},
	}

	var buf bytes.Buffer
	n.DisplayStats(&buf) // must not panic or try to allocate an enormous string
}

func TestDisplayNormalCase(t *testing.T) {
	now := time.Now()
	n := &NetworkMonitor{
		requests: map[string]*req{
			"a": {name: "first", from: now, to: now.Add(100 * time.Millisecond)},
			"b": {name: "second", from: now.Add(50 * time.Millisecond), to: now.Add(300 * time.Millisecond)},
		},
	}

	var buf bytes.Buffer
	n.DisplayStats(&buf)

	out := buf.String()
	for _, want := range []string{"first", "second", "2 requests sent"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got %q", want, out)
		}
	}
}
