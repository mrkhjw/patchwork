package validator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/validator"
)

// initRepo creates a minimal git repository and returns its path.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmds := [][]string{
		{"git", "init", dir},
		{"git", "-C", dir, "config", "user.email", "test@test.com"},
		{"git", "-C", dir, "config", "user.name", "Test"},
		{"git", "-C", dir, "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			t.Fatalf("setup cmd %v: %v", args, err)
		}
	}
	return dir
}

func writePatchFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.patch")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

func TestValidate_AllValid(t *testing.T) {
	repoDir := initRepo(t)
	patchFile := writePatchFile(t)

	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "fix-a", File: patchFile, Repos: []string{repoDir}},
		},
	}

	results := validator.Validate(cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].OK() {
		t.Errorf("expected no errors, got %v", results[0].Errors)
	}
}

func TestValidate_MissingRepo(t *testing.T) {
	patchFile := writePatchFile(t)
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "fix-b", File: patchFile, Repos: []string{"/nonexistent/repo"}},
		},
	}

	results := validator.Validate(cfg)
	if results[0].OK() {
		t.Error("expected error for missing repo")
	}
}

func TestValidate_MissingPatchFile(t *testing.T) {
	repoDir := initRepo(t)
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "fix-c", File: filepath.Join(t.TempDir(), "no.patch"), Repos: []string{repoDir}},
		},
	}

	results := validator.Validate(cfg)
	if results[0].OK() {
		t.Error("expected error for missing patch file")
	}
}

func TestAnyFailed(t *testing.T) {
	ok := validator.Result{PatchName: "a", RepoPath: "/x", Errors: nil}
	bad := validator.Result{PatchName: "b", RepoPath: "/y", Errors: []string{"oops"}}

	if validator.AnyFailed([]validator.Result{ok}) {
		t.Error("expected no failure")
	}
	if !validator.AnyFailed([]validator.Result{ok, bad}) {
		t.Error("expected failure detected")
	}
}
