package awsservices

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"

	"github.com/bootswithdefer/awless/logger"
)

func TestFileCacheProvider(t *testing.T) {
	name, err := os.MkdirTemp(".", "cache")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(name)
	os.Setenv("__AWLESS_CACHE", name)

	cacheDuration := 30 * time.Millisecond
	stsCacheDuration = cacheDuration // Force cached credential expiration after 30 milliseconds

	mock := &mockCredWithExpirationProvider{value: aws.Credentials{SecretAccessKey: "my valid secret string", Source: stscreds.ProviderName}}

	provider := fileCacheProvider{creds: mock, profile: "default", log: logger.DiscardLogger}
	ctx := context.Background()
	retrievedCredential, err := provider.Retrieve(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	retrievedCredential, err = provider.Retrieve(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	time.Sleep(cacheDuration)

	retrievedCredential, err = provider.Retrieve(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := retrievedCredential.SecretAccessKey, "my valid secret string"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := mock.accessCount, 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

type mockCredWithExpirationProvider struct {
	accessCount int
	value       aws.Credentials
}

func (m *mockCredWithExpirationProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	m.accessCount++
	return m.value, nil
}
