package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork/internal/resolver"
)

func writeTempPatch(t *testing.T, dir, name string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("--- a\n+++ b\n"), 0o644); err != nil {
		t.Fatalf("writeTempPatch: %v", err)
	}
	return p
}

func TestResolve_RelativePaths(t *testing.T) {
	dir := t.TempDir()
	writeTempPatch(t, dir, "fix.patch")

	results, err := resolver.Resolve(dir, map[string]string{
		"fix": "fix.patch",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Exists {
		t.Errorf("expected file to exist at %s", results[0].ResolvedPath)
	}
}

func TestResolve_AbsolutePath(t *testing.T) {
	dir := t.TempDir()
	abs := writeTempPatch(t, dir, "abs.patch")

	results, err := resolver.Resolve("/some/other/dir", map[string]string{
		"abs": abs,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Exists {
		t.Errorf("expected absolute path to resolve correctly")
	}
}

func TestResolve_MissingFile(t *testing.T) {
	dir := t.TempDir()

	results, err := resolver.Resolve(dir, map[string]string{
		"ghost": "nonexistent.patch",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Exists {
		t.Errorf("expected file to be reported missing")
	}
}

func TestResolve_EmptyMap(t *testing.T) {
	dir := t.TempDir()

	results, err := resolver.Resolve(dir, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty map, got %d", len(results))
	}
}

func TestAnyMissing_True(t *testing.T) {
	dir := t.TempDir()
	writeTempPatch(t, dir, "present.patch")

	results, _ := resolver.Resolve(dir, map[string]string{
		"present": "present.patch",
		"absent":  "absent.patch",
	})
	if !resolver.AnyMissing(results) {
		t.Error("expected AnyMissing to return true")
	}
}

func TestAnyMissing_False(t *testing.T) {
	dir := t.TempDir()
	writeTempPatch(t, dir, "a.patch")
	writeTempPatch(t, dir, "b.patch")

	results, _ := resolver.Resolve(dir, map[string]string{
		"a": "a.patch",
		"b": "b.patch",
	})
	if resolver.AnyMissing(results) {
		t.Error("expected AnyMissing to return false")
	}
}

func TestMissing_ReturnsOnlyAbsent(t *testing.T) {
	dir := t.TempDir()
	writeTempPatch(t, dir, "here.patch")

	results, _ := resolver.Resolve(dir, map[string]string{
		"here":  "here.patch",
		"there": "there.patch",
	})
	missing := resolver.Missing(results)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(missing))
	}
	if missing[0].PatchName != "there" {
		t.Errorf("expected missing patch name 'there', got %s", missing[0].PatchName)
	}
}
