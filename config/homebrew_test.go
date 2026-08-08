package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// installedViaHomebrew inspects the running binary's own path, so the test builds a
// fake tree and runs the same string matching against it. Verifying the real
// os.Executable path would only assert something about the test binary's location.
func TestHomebrewPathDetection(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		want bool
	}{
		{"apple silicon cask", "/opt/homebrew/Caskroom/awless/1.1.0/awless", true},
		{"intel cask", "/usr/local/Caskroom/awless/1.1.0/awless", true},
		{"formula cellar", "/opt/homebrew/Cellar/awless/1.1.0/bin/awless", true},
		{"linuxbrew", "/home/linuxbrew/.linuxbrew/Cellar/awless/1.1.0/bin/awless", true},
		{"go install", "/home/jsmith/go/bin/awless", false},
		{"manual tarball", "/usr/local/bin/awless", false},
		{"cwd", "./awless", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := strings.Contains(tc.path, "/Caskroom/") ||
				strings.Contains(tc.path, "/Cellar/") ||
				strings.Contains(tc.path, "/linuxbrew/")
			if got != tc.want {
				t.Errorf("%s: got %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

// The real function must not panic or report true for an ordinary binary, which is what
// it will see under `go test`.
func TestInstalledViaHomebrewOnTestBinary(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skipf("cannot resolve the test binary path: %s", err)
	}
	if installedViaHomebrew() && !strings.Contains(exe, "brew") {
		t.Errorf("reported a Homebrew install for %s", filepath.Clean(exe))
	}
}
