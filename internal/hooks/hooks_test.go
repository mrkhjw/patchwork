package hooks_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/patchwork/internal/hooks"
)

func TestRunAll_Success(t *testing.T) {
	dir := t.TempDir()
	h := []hooks.Hook{{Command: "echo hello", Dir: dir}}
	results := hooks.RunAll(h, dir)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
	if results[0].Output != "hello" {
		t.Errorf("expected output 'hello', got %q", results[0].Output)
	}
}

func TestRunAll_StopsOnFailure(t *testing.T) {
	dir := t.TempDir()
	h := []hooks.Hook{
		{Command: "false"},
		{Command: "echo should-not-run"},
	}
	results := hooks.RunAll(h, dir)
	if len(results) != 1 {
		t.Fatalf("expected early stop after 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected error from 'false', got nil")
	}
}

func TestRunAll_EmptyCommand(t *testing.T) {
	dir := t.TempDir()
	h := []hooks.Hook{{Command: ""}}
	results := hooks.RunAll(h, dir)
	if results[0].Err == nil {
		t.Error("expected error for empty command")
	}
}

func TestRunAll_UsesDefaultDir(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "marker.txt")
	if err := os.WriteFile(marker, []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	// Hook with no Dir set — should fall back to defaultDir
	h := []hooks.Hook{{Command: "ls marker.txt"}}
	results := hooks.RunAll(h, dir)
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := hooks.RunAll([]hooks.Hook{{Command: "false"}}, t.TempDir())
	if !hooks.AnyFailed(results) {
		t.Error("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := hooks.RunAll([]hooks.Hook{{Command: "echo ok"}}, t.TempDir())
	if hooks.AnyFailed(results) {
		t.Error("expected AnyFailed to return false")
	}
}
