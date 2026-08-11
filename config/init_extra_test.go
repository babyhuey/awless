package config

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bootswithdefer/awless/database"
)

// withAwlessHome redirects every path the config package derives from $HOME into a temp
// directory. These are package-level vars computed once at init, so they are the seam:
// without overriding them a test would create ~/.awless on the developer's machine.
func withAwlessHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	prevHome, prevDB, prevDir, prevKeys := AwlessHome, DBPath, Dir, KeysDir
	prevFirst := AwlessFirstInstall
	AwlessHome = dir
	DBPath = filepath.Join(dir, database.Filename)
	Dir = filepath.Join(dir, "aws")
	KeysDir = filepath.Join(dir, "keys")

	// database.Execute resolves its file from this env var rather than from DBPath.
	t.Setenv("__AWLESS_HOME", dir)

	// Without this, region resolution falls through to the EC2 instance metadata
	// service and waits out its timeout on any machine that is not an EC2 instance.
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	// InitAwlessEnv writes an initial config on first install, and that write rejects
	// an empty region.
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	t.Cleanup(func() {
		AwlessHome, DBPath, Dir, KeysDir = prevHome, prevDB, prevDir, prevKeys
		AwlessFirstInstall = prevFirst
	})
	return dir
}

// The keys directory holds private SSH keys, so its mode is security-relevant: 0700 and
// nothing wider.
func TestInitAwlessEnvCreatesKeysDirPrivately(t *testing.T) {
	dir := withAwlessHome(t)

	if err := InitAwlessEnv(); err != nil {
		t.Fatalf("InitAwlessEnv: %s", err)
	}

	info, err := os.Stat(filepath.Join(dir, "keys"))
	if err != nil {
		t.Fatalf("keys dir was not created: %s", err)
	}
	if !info.IsDir() {
		t.Fatal("keys path is not a directory")
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("keys dir mode = %#o, want 0700", perm)
	}
}

// AwlessFirstInstall drives the welcome banner and the initial config write, and is
// exported through an env var that other packages read.
func TestInitAwlessEnvDetectsFirstInstall(t *testing.T) {
	dir := withAwlessHome(t)

	if err := InitAwlessEnv(); err != nil {
		t.Fatalf("InitAwlessEnv: %s", err)
	}
	if !AwlessFirstInstall {
		t.Error("a fresh directory should be reported as a first install")
	}
	if got := os.Getenv("__AWLESS_FIRST_INSTALL"); got != "true" {
		t.Errorf("__AWLESS_FIRST_INSTALL = %q, want true", got)
	}

	// Second run against the same home, where the database now exists.
	if _, err := os.Stat(filepath.Join(dir, database.Filename)); err != nil {
		t.Skipf("no database was created, cannot exercise the second run: %s", err)
	}
	if err := InitAwlessEnv(); err != nil {
		t.Fatalf("second InitAwlessEnv: %s", err)
	}
	if AwlessFirstInstall {
		t.Error("an existing database should not be reported as a first install")
	}
	if got := os.Getenv("__AWLESS_FIRST_INSTALL"); got != "false" {
		t.Errorf("__AWLESS_FIRST_INSTALL = %q, want false", got)
	}
}

func TestResolveRequiredConfigFromEnv(t *testing.T) {
	t.Run("region present", func(t *testing.T) {
		t.Setenv("AWS_DEFAULT_REGION", "eu-west-3")
		t.Setenv("AWS_REGION", "")
		t.Setenv("AWS_EC2_METADATA_DISABLED", "true")

		got := resolveRequiredConfigFromEnv()
		if got[RegionConfigKey] != "eu-west-3" {
			t.Errorf("resolved region = %q, want eu-west-3", got[RegionConfigKey])
		}
	})

	t.Run("no region resolves to an empty map", func(t *testing.T) {
		t.Setenv("AWS_DEFAULT_REGION", "")
		t.Setenv("AWS_REGION", "")
		t.Setenv("AWS_EC2_METADATA_DISABLED", "true")

		got := resolveRequiredConfigFromEnv()
		// An empty string must not be recorded as the region: it would be written to
		// config as a valid-looking value.
		if v, ok := got[RegionConfigKey]; ok && v == "" {
			t.Error("an empty region was recorded in the resolved config")
		}
	})
}

// VerifyNewVersionAvailable gates the network call on a stored timestamp. The gating is
// the part worth pinning: an always-on check would hit GitHub on every command.
func TestVerifyNewVersionAvailableRespectsFrequency(t *testing.T) {
	withAwlessHome(t)

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		w.Write([]byte(`{"Version":"v99.0.0","URL":"https://example.invalid"}`)) //nolint:errcheck // test server
	}))
	t.Cleanup(srv.Close)

	var messaging bytes.Buffer

	// No timestamp stored yet, so the check is due and should run.
	if err := VerifyNewVersionAvailable(context.Background(), srv.URL, &messaging); err != nil {
		t.Fatalf("first check: %s", err)
	}
	if hits != 1 {
		t.Fatalf("server hit %d times on a due check, want 1", hits)
	}
	if messaging.Len() == 0 {
		t.Error("an available upgrade produced no message")
	}

	// The timestamp is now fresh, so an immediate second call must not call out again.
	if err := VerifyNewVersionAvailable(context.Background(), srv.URL, &messaging); err != nil {
		t.Fatalf("second check: %s", err)
	}
	if hits != 1 {
		t.Errorf("server hit %d times, want 1 — the frequency gate did not hold", hits)
	}

	// Confirm the gate is the stored timestamp rather than a process-lifetime flag.
	err := database.Execute(func(db *database.DB) error {
		return db.SetTimeValue(lastUpgradeCheckDBKey, time.Now().Add(-365*24*time.Hour))
	})
	if err != nil {
		t.Fatalf("backdating the last check: %s", err)
	}
	if err := VerifyNewVersionAvailable(context.Background(), srv.URL, &messaging); err != nil {
		t.Fatalf("third check: %s", err)
	}
	if hits != 2 {
		t.Errorf("server hit %d times after backdating, want 2", hits)
	}
}

// A negative frequency disables the check entirely, which is how a user opts out.
func TestVerifyNewVersionAvailableCanBeDisabled(t *testing.T) {
	withAwlessHome(t)

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { hits++ }))
	t.Cleanup(srv.Close)

	if err := database.Execute(func(db *database.DB) error {
		return db.SetTimeValue(lastUpgradeCheckDBKey, time.Time{})
	}); err != nil {
		t.Fatalf("clearing the last check: %s", err)
	}

	// The frequency is read from the in-memory Config map rather than straight from
	// the database, so that is what has to be negative here.
	prev, had := Config[checkUpgradeFrequencyConfigKey]
	Config[checkUpgradeFrequencyConfigKey] = -1
	t.Cleanup(func() {
		if had {
			Config[checkUpgradeFrequencyConfigKey] = prev
		} else {
			delete(Config, checkUpgradeFrequencyConfigKey)
		}
	})

	var messaging bytes.Buffer
	if err := VerifyNewVersionAvailable(context.Background(), srv.URL, &messaging); err != nil {
		t.Fatalf("check: %s", err)
	}
	if hits != 0 {
		t.Errorf("server hit %d times while disabled, want 0", hits)
	}
}
