/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package repo

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func TestReduceToLastRevOfEachDay(t *testing.T) {
	revs := []*Rev{
		{ID: "1", Date: mustParse("2017-01-18 15:05")},
		{ID: "2", Date: mustParse("2017-01-18 15:09")},
		{ID: "3", Date: mustParse("2017-01-19 09:05")},
		{ID: "4", Date: mustParse("2017-01-19 08:05")},
		{ID: "5", Date: mustParse("2017-01-17 21:05")},
		{ID: "6", Date: mustParse("2017-01-17 10:05")},
	}

	reduced := reduceToLastRevOfEachDay(revs)

	sort.Slice(reduced, func(i, j int) bool { return reduced[i].Date.Before(reduced[j].Date) })

	if got, want := len(reduced), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := reduced[0].ID, "5"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := reduced[1].ID, "2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := reduced[2].ID, "3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestReduceToLastRevOfEachDayEmpty(t *testing.T) {
	reduced := reduceToLastRevOfEachDay(nil)
	if got := len(reduced); got != 0 {
		t.Fatalf("expected 0 revs, got %d", got)
	}
}

func TestReduceToLastRevOfEachDaySingleDay(t *testing.T) {
	revs := []*Rev{
		{ID: "a", Date: mustParse("2017-03-01 08:00")},
		{ID: "b", Date: mustParse("2017-03-01 12:00")},
		{ID: "c", Date: mustParse("2017-03-01 18:00")},
	}

	reduced := reduceToLastRevOfEachDay(revs)

	if got := len(reduced); got != 1 {
		t.Fatalf("expected 1 rev, got %d", got)
	}
	if got := reduced[0].ID; got != "c" {
		t.Fatalf("expected latest rev 'c', got %s", got)
	}
}

func TestReduceToLastRevOfEachDaySingleRev(t *testing.T) {
	revs := []*Rev{
		{ID: "only", Date: mustParse("2020-06-15 10:30")},
	}
	reduced := reduceToLastRevOfEachDay(revs)
	if got := len(reduced); got != 1 {
		t.Fatalf("expected 1 rev, got %d", got)
	}
	if got := reduced[0].ID; got != "only" {
		t.Fatalf("expected rev 'only', got %s", got)
	}
}

func TestRevDateString(t *testing.T) {
	rev := &Rev{
		ID:   "abc123",
		Date: time.Date(2017, 1, 18, 15, 5, 30, 0, time.UTC),
	}
	got := rev.DateString()
	exp := "Wed Jan 18 15:05:30"
	if got != exp {
		t.Fatalf("got %q, want %q", got, exp)
	}
}

func TestRevDateStringDifferentDate(t *testing.T) {
	rev := &Rev{
		ID:   "def456",
		Date: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
	}
	got := rev.DateString()
	exp := "Mon Dec 25 00:00:00"
	if got != exp {
		t.Fatalf("got %q, want %q", got, exp)
	}
}

func TestNullRepoCommit(t *testing.T) {
	nr := NullRepo{}
	if err := nr.Commit("file1", "file2"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNullRepoList(t *testing.T) {
	nr := NullRepo{}
	revs, err := nr.List()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if revs != nil {
		t.Fatalf("expected nil revs, got %v", revs)
	}
}

func TestNullRepoLoadRev(t *testing.T) {
	nr := NullRepo{}
	rev, err := nr.LoadRev("some-version")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if rev != nil {
		t.Fatalf("expected nil rev, got %v", rev)
	}
}

func TestNullRepoBaseDir(t *testing.T) {
	nr := NullRepo{}
	if got := nr.BaseDir(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestBaseDir(t *testing.T) {
	orig := os.Getenv("__AWLESS_HOME")
	defer os.Setenv("__AWLESS_HOME", orig)

	os.Setenv("__AWLESS_HOME", "/tmp/test-awless-home")
	expected := filepath.Join("/tmp/test-awless-home", "aws", "rdf")
	if got := BaseDir(); got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestBaseDirEmpty(t *testing.T) {
	orig := os.Getenv("__AWLESS_HOME")
	defer os.Setenv("__AWLESS_HOME", orig)

	os.Setenv("__AWLESS_HOME", "")
	expected := filepath.Join("", "aws", "rdf")
	if got := BaseDir(); got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestNewGitRepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-repo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatalf("unexpected error creating git repo: %v", err)
	}
	if repo == nil {
		t.Fatal("expected non-nil repo")
	}
	if got := repo.BaseDir(); got != dir {
		t.Fatalf("expected basedir %q, got %q", dir, got)
	}

	// Verify .git directory was created
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Fatal("expected .git directory to be created")
	}
}

func TestNewGitRepoAlreadyInitialized(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-repo-existing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Init first time
	repo1, err := newGitRepo(dir)
	if err != nil {
		t.Fatalf("unexpected error on first init: %v", err)
	}

	// Init second time - should reuse existing
	repo2, err := newGitRepo(dir)
	if err != nil {
		t.Fatalf("unexpected error on second init: %v", err)
	}
	if repo1.BaseDir() != repo2.BaseDir() {
		t.Fatalf("basedirs should match: %q vs %q", repo1.BaseDir(), repo2.BaseDir())
	}
}

func TestGitRepoCommitAndList(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-commit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a file to commit
	testFile := "test.txt"
	if err := os.WriteFile(filepath.Join(dir, testFile), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.Commit(testFile); err != nil {
		t.Fatalf("unexpected error committing: %v", err)
	}

	// List revisions
	revs, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(revs) != 1 {
		t.Fatalf("expected 1 revision, got %d", len(revs))
	}
	if revs[0].ID == "" {
		t.Fatal("expected non-empty revision Id")
	}
	if revs[0].Date.IsZero() {
		t.Fatal("expected non-zero revision Date")
	}
}

func TestGitRepoMultipleCommits(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-multi-commit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatal(err)
	}

	// First commit
	if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := repo.Commit("file1.txt"); err != nil {
		t.Fatal(err)
	}

	// Second commit
	if err := os.WriteFile(filepath.Join(dir, "file2.txt"), []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := repo.Commit("file2.txt"); err != nil {
		t.Fatal(err)
	}

	revs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 2 {
		t.Fatalf("expected 2 revisions, got %d", len(revs))
	}

	// Should be sorted by date
	if !revs[0].Date.Before(revs[1].Date) && !revs[0].Date.Equal(revs[1].Date) {
		t.Fatal("expected revisions to be sorted by date ascending")
	}
}

func TestGitRepoLoadRev(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-loadrev")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Create infra and access triple files (can be empty or minimal)
	infraContent := ""
	accessContent := ""
	if err := os.WriteFile(filepath.Join(dir, "infra.triples"), []byte(infraContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "access.triples"), []byte(accessContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.Commit("infra.triples", "access.triples"); err != nil {
		t.Fatal(err)
	}

	// Get the revision
	revs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 1 {
		t.Fatalf("expected 1 rev, got %d", len(revs))
	}

	// Load the revision
	rev, err := repo.LoadRev(revs[0].ID)
	if err != nil {
		t.Fatalf("unexpected error loading rev: %v", err)
	}
	if rev == nil {
		t.Fatal("expected non-nil rev")
	}
	if rev.ID != revs[0].ID {
		t.Fatalf("expected Id %s, got %s", revs[0].ID, rev.ID)
	}
	if rev.Infra == nil {
		t.Fatal("expected non-nil Infra graph")
	}
	if rev.Access == nil {
		t.Fatal("expected non-nil Access graph")
	}
	if rev.Date.IsZero() {
		t.Fatal("expected non-zero Date")
	}
}

func TestGitRepoLoadRevInvalidHash(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-loadrev-invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a commit first so repo is not empty
	if err := os.WriteFile(filepath.Join(dir, "dummy.txt"), []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := repo.Commit("dummy.txt"); err != nil {
		t.Fatal(err)
	}

	// Try to load with invalid hash
	_, err = repo.LoadRev("0000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected error for invalid hash")
	}
}

func TestGitRepoListEmpty(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-git-list-empty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repo, err := newGitRepo(dir)
	if err != nil {
		t.Fatal(err)
	}

	revs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 0 {
		t.Fatalf("expected 0 revisions for empty repo, got %d", len(revs))
	}
}

func TestNewRepoSetsEnv(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-repo-new")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	orig := os.Getenv("__AWLESS_HOME")
	os.Setenv("__AWLESS_HOME", tmpDir)
	defer os.Setenv("__AWLESS_HOME", orig)

	repo, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(tmpDir, "aws", "rdf")
	if got := repo.BaseDir(); got != expected {
		t.Fatalf("expected basedir %q, got %q", expected, got)
	}
}

func mustParse(s string) time.Time {
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return t
}
