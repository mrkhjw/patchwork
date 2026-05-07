package linter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/linter"
)

func writeTempPatch(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "fix.patch")
	if err := os.WriteFile(p, []byte("--- a\n+++ b\n"), 0o644); err != nil {
		t.Fatalf("write patch: %v", err)
	}
	return p
}

func TestLint_ValidPatch(t *testing.T) {
	patchFile := writeTempPatch(t)
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "my-fix", File: patchFile, Repo: "/some/repo", Tags: []string{"hotfix"}},
		},
	}
	results := linter.Lint(cfg)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].AnyErrors() {
		t.Errorf("unexpected errors: %v", results[0].Errors)
	}
	if results[0].AnyWarnings() {
		t.Errorf("unexpected warnings: %v", results[0].Warnings)
	}
}

func TestLint_MissingPatchFile(t *testing.T) {
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "bad", File: "/nonexistent/fix.patch", Repo: "/repo", Tags: []string{"x"}},
		},
	}
	results := linter.Lint(cfg)
	if !results[0].AnyErrors() {
		t.Fatal("expected errors for missing patch file")
	}
}

func TestLint_EmptyName(t *testing.T) {
	patchFile := writeTempPatch(t)
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "", File: patchFile, Repo: "/repo", Tags: []string{"x"}},
		},
	}
	results := linter.Lint(cfg)
	if !results[0].AnyErrors() {
		t.Fatal("expected error for empty name")
	}
}

func TestLint_NoTags_ProducesWarning(t *testing.T) {
	patchFile := writeTempPatch(t)
	cfg := &config.Config{
		Patches: []config.Patch{
			{Name: "notag", File: patchFile, Repo: "/repo", Tags: nil},
		},
	}
	results := linter.Lint(cfg)
	if results[0].AnyErrors() {
		t.Fatalf("unexpected errors: %v", results[0].Errors)
	}
	if !results[0].AnyWarnings() {
		t.Fatal("expected warning for missing tags")
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := []linter.Result{
		{PatchName: "ok"},
		{PatchName: "bad", Errors: []string{"something wrong"}},
	}
	if !linter.AnyFailed(results) {
		t.Fatal("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := []linter.Result{
		{PatchName: "ok"},
	}
	if linter.AnyFailed(results) {
		t.Fatal("expected AnyFailed to return false")
	}
}
