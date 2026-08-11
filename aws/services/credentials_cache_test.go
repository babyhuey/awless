package awsservices

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"

	"github.com/bootswithdefer/awless/logger"
)

// The credential's own expiry is authoritative. A profile with duration_seconds=3600 was
// previously re-prompting for MFA four times an hour, because the cache assumed a flat
// fifteen minutes regardless of what the credential said.
func TestCacheExpiryPrefersTheCredentialExpiration(t *testing.T) {
	expires := time.Now().UTC().Add(1 * time.Hour)
	got := cacheExpiry(aws.Credentials{CanExpire: true, Expires: expires})

	want := expires.Add(-expiryMargin)
	if !got.Equal(want) {
		t.Errorf("got %s, want %s (the credential expiry less the safety margin)", got, want)
	}
	// The flat fallback would land ~45 minutes earlier; make sure that is not what
	// happened.
	if got.Before(time.Now().UTC().Add(30 * time.Minute)) {
		t.Errorf("expiry %s ignored the credential's own one-hour lifetime", got)
	}
}

// A margin keeps a long sync from starting with a valid credential and finishing without.
func TestCacheExpiryKeepsASafetyMargin(t *testing.T) {
	expires := time.Now().UTC().Add(1 * time.Hour)
	got := cacheExpiry(aws.Credentials{CanExpire: true, Expires: expires})

	if !got.Before(expires) {
		t.Errorf("expiry %s should be before the credential's own %s", got, expires)
	}
}

// Some providers report no expiry at all, and those fall back to the flat duration.
func TestCacheExpiryFallsBackWhenTheCredentialDoesNotExpire(t *testing.T) {
	for _, c := range []aws.Credentials{
		{CanExpire: false},
		{CanExpire: true}, // CanExpire, but a zero Expires
	} {
		got := cacheExpiry(c)
		lower := time.Now().UTC().Add(stsCacheDuration - time.Minute)
		if got.Before(lower) {
			t.Errorf("CanExpire=%v: got %s, want roughly %s from now", c.CanExpire, got, stsCacheDuration)
		}
	}
}

// Static keys already sit in the shared credentials file. Copying them into the awless
// cache would widen their exposure and buy nothing, so only assumed-role credentials are
// written out.
func TestFileCacheProviderDoesNotWriteStaticCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("__AWLESS_CACHE", dir)

	mock := &mockCredWithExpirationProvider{value: aws.Credentials{
		SecretAccessKey: "static-secret",
		Source:          "SharedConfigCredentials", // not stscreds.ProviderName
	}}
	provider := &fileCacheProvider{creds: mock, profile: "default", log: logger.DiscardLogger}

	if _, err := provider.Retrieve(context.Background()); err != nil {
		t.Fatal(err)
	}

	credFile := filepath.Join(dir, "credentials", "aws-profile-default.json")
	if _, err := os.Stat(credFile); !os.IsNotExist(err) {
		t.Errorf("static credentials were written to %s; only assumed-role ones should be cached", credFile)
	}
}

func TestFileCacheProviderWritesAssumedRoleCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("__AWLESS_CACHE", dir)

	mock := &mockCredWithExpirationProvider{value: aws.Credentials{
		SecretAccessKey: "assumed-secret",
		Source:          stscreds.ProviderName,
		CanExpire:       true,
		Expires:         time.Now().UTC().Add(1 * time.Hour),
	}}
	provider := &fileCacheProvider{creds: mock, profile: "myprofile", log: logger.DiscardLogger}

	if _, err := provider.Retrieve(context.Background()); err != nil {
		t.Fatal(err)
	}

	credFile := filepath.Join(dir, "credentials", "aws-profile-myprofile.json")
	info, err := os.Stat(credFile)
	if err != nil {
		t.Fatalf("expected a cache file at %s: %s", credFile, err)
	}
	// The file holds live credentials, so it must not be group- or world-readable.
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("cache file mode is %04o, want 0600", perm)
	}
}

// An empty profile means the default one. Without normalizing it the file was named
// "aws-profile-.json".
func TestFileCacheProviderNamesTheDefaultProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("__AWLESS_CACHE", dir)

	mock := &mockCredWithExpirationProvider{value: aws.Credentials{
		SecretAccessKey: "assumed-secret",
		Source:          stscreds.ProviderName,
	}}
	provider := &fileCacheProvider{creds: mock, profile: "", log: logger.DiscardLogger}

	if _, err := provider.Retrieve(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "credentials", "aws-profile-default.json")); err != nil {
		t.Errorf("expected the default profile cache file: %s", err)
	}
}

// A second process reads what the first wrote, which is the whole point: MFA once per
// session rather than once per invocation.
func TestFileCacheProviderReadsAnotherProcessesCache(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("__AWLESS_CACHE", dir)

	first := &mockCredWithExpirationProvider{value: aws.Credentials{
		SecretAccessKey: "assumed-secret",
		Source:          stscreds.ProviderName,
		CanExpire:       true,
		Expires:         time.Now().UTC().Add(1 * time.Hour),
	}}
	if _, err := (&fileCacheProvider{creds: first, profile: "p", log: logger.DiscardLogger}).
		Retrieve(context.Background()); err != nil {
		t.Fatal(err)
	}

	// A fresh provider, as a new invocation would build: no in-memory state, so a cache
	// hit can only come from the file.
	second := &mockCredWithExpirationProvider{value: aws.Credentials{SecretAccessKey: "should-not-be-used"}}
	got, err := (&fileCacheProvider{creds: second, profile: "p", log: logger.DiscardLogger}).
		Retrieve(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if got.SecretAccessKey != "assumed-secret" {
		t.Errorf("got %q, want the cached credential", got.SecretAccessKey)
	}
	if second.accessCount != 0 {
		t.Errorf("the underlying provider was called %d times; the cache should have served this", second.accessCount)
	}
}

// Serving from memory is what keeps the file from being stat'd and read on every single
// AWS request.
func TestFileCacheProviderServesRepeatCallsFromMemory(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("__AWLESS_CACHE", dir)

	mock := &mockCredWithExpirationProvider{value: aws.Credentials{
		SecretAccessKey: "assumed-secret",
		Source:          stscreds.ProviderName,
		CanExpire:       true,
		Expires:         time.Now().UTC().Add(1 * time.Hour),
	}}
	provider := &fileCacheProvider{creds: mock, profile: "p", log: logger.DiscardLogger}
	ctx := context.Background()

	if _, err := provider.Retrieve(ctx); err != nil {
		t.Fatal(err)
	}
	// Remove the file: if a later call still succeeds, it was served from memory.
	if err := os.RemoveAll(filepath.Join(dir, "credentials")); err != nil {
		t.Fatal(err)
	}

	for i := range 5 {
		got, err := provider.Retrieve(ctx)
		if err != nil {
			t.Fatalf("call %d: %s", i, err)
		}
		if got.SecretAccessKey != "assumed-secret" {
			t.Errorf("call %d: got %q", i, got.SecretAccessKey)
		}
	}
	if mock.accessCount != 1 {
		t.Errorf("underlying provider called %d times, want 1", mock.accessCount)
	}
}
