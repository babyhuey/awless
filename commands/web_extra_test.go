package commands

import (
	"net"
	"strings"
	"testing"
)

// The web UI serves the local synced AWS inventory with no authentication, so
// it must not bind a routable interface unless explicitly asked to.
func TestWebListenAddress(t *testing.T) {
	tcases := []struct {
		desc      string
		port      string
		listenAll bool
		want      string
	}{
		{desc: "default is loopback", port: ":8080", listenAll: false, want: "127.0.0.1:8080"},
		{desc: "port without colon", port: "8080", listenAll: false, want: "127.0.0.1:8080"},
		{desc: "custom port stays loopback", port: ":9090", listenAll: false, want: "127.0.0.1:9090"},
		{desc: "listen-all opts in explicitly", port: ":8080", listenAll: true, want: "0.0.0.0:8080"},
	}

	for _, tc := range tcases {
		t.Run(tc.desc, func(t *testing.T) {
			got := webListenAddr(tc.port, tc.listenAll)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
			if _, _, err := net.SplitHostPort(got); err != nil {
				t.Errorf("%q is not a valid listen address: %s", got, err)
			}
		})
	}
}

func TestWebDefaultNeverBindsWildcard(t *testing.T) {
	for _, port := range []string{":8080", "8080", ":1", ":65535"} {
		addr := webListenAddr(port, false)
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			t.Fatalf("invalid addr %q: %s", addr, err)
		}
		if !net.ParseIP(host).IsLoopback() {
			t.Errorf("default bind host %q for port %q is not loopback", host, port)
		}
		if strings.HasPrefix(addr, ":") || strings.HasPrefix(addr, "0.0.0.0") {
			t.Errorf("default addr %q binds all interfaces", addr)
		}
	}
}
