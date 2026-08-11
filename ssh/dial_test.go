package ssh

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	gossh "golang.org/x/crypto/ssh"

	"github.com/bootswithdefer/awless/logger"
)

// newDialTestClient builds a Client wired to a fake dialer. The logger is set
// explicitly because DialWithUsers logs on success and a zero-value Client has a nil
// *logger.Logger, whose methods dereference the receiver.
func newDialTestClient(t *testing.T, dial func(string, string, *gossh.ClientConfig) (*gossh.Client, error)) *Client {
	t.Helper()
	return &Client{
		Config:                &gossh.ClientConfig{User: "original"},
		IP:                    "10.0.0.1",
		Port:                  22,
		StrictHostKeyChecking: true,
		logger:                logger.DiscardLogger,
		dialFunc:              dial,
	}
}

// DialWithUsers tries candidate usernames in order and keeps the first that works.
// This is the behavior awless relies on to guess the login for an AMI whose default
// user it does not know, so both the stopping point and the recorded user matter.
func TestDialWithUsersStopsAtFirstSuccess(t *testing.T) {
	var attempted []string
	c := newDialTestClient(t, func(_, _ string, config *gossh.ClientConfig) (*gossh.Client, error) {
		attempted = append(attempted, config.User)
		if config.User == "ec2-user" {
			return &gossh.Client{}, nil
		}
		return nil, errors.New("permission denied")
	})

	if err := c.DialWithUsers("root", "admin", "ec2-user", "ubuntu"); err != nil {
		t.Fatalf("expected success, got %s", err)
	}
	if c.User != "ec2-user" {
		t.Errorf("recorded user = %q, want ec2-user", c.User)
	}
	if c.Client == nil {
		t.Error("the successful client was not retained")
	}
	// "ubuntu" must not be tried: a further dial after success is a wasted
	// connection attempt against a host that has already authenticated.
	want := []string{"root", "admin", "ec2-user"}
	if strings.Join(attempted, ",") != strings.Join(want, ",") {
		t.Errorf("attempted %v, want %v", attempted, want)
	}
}

func TestDialWithUsersTriesEveryUserBeforeFailing(t *testing.T) {
	lastErr := errors.New("no supported methods remain")
	var attempted []string
	c := newDialTestClient(t, func(_, _ string, config *gossh.ClientConfig) (*gossh.Client, error) {
		attempted = append(attempted, config.User)
		return nil, lastErr
	})

	err := c.DialWithUsers("root", "admin", "ubuntu")
	if err == nil {
		t.Fatal("expected an error when every username fails")
	}
	if len(attempted) != 3 {
		t.Errorf("attempted %d users, want 3", len(attempted))
	}
	// The last error is wrapped so callers can inspect the underlying cause.
	if !errors.Is(err, lastErr) {
		t.Errorf("error does not wrap the last dial error: %v", err)
	}
	// The message must name the host and the users tried, since that is the only
	// clue the user gets about which logins were attempted.
	for _, want := range []string{"10.0.0.1:22", "root", "admin", "ubuntu"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
	if c.User != "" || c.Client != nil {
		t.Error("a failed dial must not record a user or client")
	}
}

// An empty username list used to fall through the loop with a nil error and wrap it,
// rendering the unhelpful "Last error: %!w(<nil>)".
func TestDialWithUsersWithNoUsernames(t *testing.T) {
	called := false
	c := newDialTestClient(t, func(string, string, *gossh.ClientConfig) (*gossh.Client, error) {
		called = true
		return &gossh.Client{}, nil
	})

	err := c.DialWithUsers()
	if err == nil {
		t.Fatal("expected an error when no username is given")
	}
	if called {
		t.Error("must not dial when there is no username to try")
	}
	if strings.Contains(err.Error(), "%!w") {
		t.Errorf("error wraps a nil error: %q", err)
	}
	if errors.Unwrap(err) != nil {
		t.Errorf("nothing should be wrapped, got %v", errors.Unwrap(err))
	}
}

// Each attempt copies the config, so the caller's ClientConfig — shared with the proxy
// client and reused across attempts — must come back unmodified.
func TestDialWithUsersDoesNotMutateSharedConfig(t *testing.T) {
	original := &gossh.ClientConfig{User: "original"}
	c := newDialTestClient(t, func(_, _ string, config *gossh.ClientConfig) (*gossh.Client, error) {
		if config == original {
			t.Error("dialed with the caller's config rather than a copy")
		}
		return nil, errors.New("nope")
	})
	c.Config = original

	_ = c.DialWithUsers("root", "admin")

	if original.User != "original" {
		t.Errorf("caller config User mutated to %q", original.User)
	}
	if original.HostKeyCallback != nil {
		t.Error("caller config HostKeyCallback was mutated")
	}
}

// Host key checking is the one security-relevant branch here: it must be bypassed only
// when explicitly disabled, and the real callback preserved otherwise.
func TestDialWithUsersHostKeyChecking(t *testing.T) {
	sentinel := errors.New("host key rejected")

	for _, tc := range []struct {
		name       string
		strict     bool
		wantStrict bool
	}{
		{"strict keeps the configured callback", true, true},
		{"lax substitutes an insecure callback", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var seen gossh.HostKeyCallback
			c := newDialTestClient(t, func(_, _ string, config *gossh.ClientConfig) (*gossh.Client, error) {
				seen = config.HostKeyCallback
				return &gossh.Client{}, nil
			})
			c.Config = &gossh.ClientConfig{
				HostKeyCallback: func(string, net.Addr, gossh.PublicKey) error { return sentinel },
			}
			c.StrictHostKeyChecking = tc.strict

			if err := c.DialWithUsers("root"); err != nil {
				t.Fatalf("dial: %s", err)
			}
			if seen == nil {
				t.Fatal("no host key callback was set")
			}
			// The configured callback rejects; the insecure one accepts. That
			// difference identifies which was used without inspecting internals.
			gotStrict := seen("host", &net.TCPAddr{}, nil) != nil
			if gotStrict != tc.wantStrict {
				t.Errorf("strict callback used = %v, want %v", gotStrict, tc.wantStrict)
			}
		})
	}
}

// The production path must not require the seam: a Client built without dialFunc has to
// fall back to the real dialler rather than nil-panic.
func TestDialerFallsBackToRealDial(t *testing.T) {
	c := &Client{}
	if c.dialer() == nil {
		t.Fatal("dialer() returned nil for a zero-value Client")
	}
	withSeam := &Client{dialFunc: func(string, string, *gossh.ClientConfig) (*gossh.Client, error) {
		return nil, fmt.Errorf("seam")
	}}
	if _, err := withSeam.dialer()("tcp", "x", nil); err == nil || err.Error() != "seam" {
		t.Errorf("dialer() did not return the injected function, got %v", err)
	}
}
