package repo_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork/internal/repo"
)

// initRepo creates a temporary git repository and returns its path.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	return dir
}

func TestExists_Valid(t *testing.T) {
	dir := initRepo(t)
	if !repo.Exists(dir) {
		t.Errorf("expected Exists to return true for %s", dir)
	}
}

func TestExists_Invalid(t *testing.T) {
	dir := t.TempDir()
	if repo.Exists(dir) {
		t.Errorf("expected Exists to return false for non-git dir")
	}
}

func TestCurrentBranch(t *testing.T) {
	dir := initRepo(t)
	// Create an initial commit so HEAD resolves to a branch.
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "."}, {"commit", "-m", "init"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	branch, err := repo.CurrentBranch(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestIsClean(t *testing.T) {
	dir := initRepo(t)
	clean, err := repo.IsClean(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !clean {
		t.Error("expected fresh repo to be clean")
	}

	// Dirty the repo.
	if err := os.WriteFile(filepath.Join(dir, "dirty.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	clean, err = repo.IsClean(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if clean {
		t.Error("expected dirty repo to not be clean")
	}
}

func TestGetStatus_NotARepo(t *testing.T) {
	dir := t.TempDir()
	_, err := repo.GetStatus(dir)
	if err == nil {
		t.Error("expected error for non-git directory")
	}
}
