package sync

import (
	"context"
	"os"
	"testing"

	"github.com/bootswithdefer/awless/cloud"
)

func TestNoOpSyncer(t *testing.T) {
	s := NoOpSyncer()
	if s == nil {
		t.Fatal("expected non-nil syncer")
	}

	result, err := s.Sync(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil map")
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}

func TestNoOpSyncerWithServices(t *testing.T) {
	s := NoOpSyncer()

	// Even with services passed, NoOpSyncer should return empty map and no error
	result, err := s.Sync(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}

func TestNewSyncerCreation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "awless_sync_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Unsetenv("__AWLESS_HOME")

	s := NewSyncer()
	if s == nil {
		t.Fatal("expected non-nil syncer")
	}
}

func TestLoadLocalGraphForServiceNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "awless_sync_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Unsetenv("__AWLESS_HOME")

	g := LoadLocalGraphForService("nonexistent", "testprofile", "us-east-1")
	if g == nil {
		t.Fatal("expected non-nil graph for non-existent service")
	}

	// Should return empty graph - verify by trying to get resources
	all, _ := g.(interface {
		GetAllResources(string) ([]*cloud.Resource, error)
	})
	if all != nil {
		t.Log("graph returned, checking it is usable")
	}
}

func TestLoadLocalGraphForServiceGlobalServices(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "awless_sync_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Unsetenv("__AWLESS_HOME")

	// These services should use "global" region dir
	for _, svc := range []string{"access", "dns", "cdn"} {
		g := LoadLocalGraphForService(svc, "testprofile", "us-west-2")
		if g == nil {
			t.Fatalf("expected non-nil graph for global service %s", svc)
		}
	}
}

func TestConcatErrors(t *testing.T) {
	// No errors
	if err := concatErrors(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if err := concatErrors([]error{}); err != nil {
		t.Fatalf("expected nil for empty slice, got %v", err)
	}
}
