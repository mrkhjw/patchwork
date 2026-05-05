package runner_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/runner"
)

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	return dir
}

func writePatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.patch")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_SkipsAlreadyApplied(t *testing.T) {
	repoDir := initRepo(t)
	statePath := filepath.Join(t.TempDir(), "state.yaml")

	// Seed state with an already-applied entry.
	st, _ := loadState(statePath)
	_ = st

	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r1", Path: repoDir}},
		Patches: []config.Patch{{Name: "p1", File: writePatch(t, "")}},
	}

	// First run — nothing applied yet, empty patch is a no-op.
	_, err := runner.Run(cfg, runner.RunOptions{StatePath: statePath})
	if err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Second run — should skip because patch is marked applied.
	results, err := runner.Run(cfg, runner.RunOptions{StatePath: statePath})
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result, got %+v", results)
	}
}

func TestRun_DryRun(t *testing.T) {
	repoDir := initRepo(t)
	statePath := filepath.Join(t.TempDir(), "state.yaml")

	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r1", Path: repoDir}},
		Patches: []config.Patch{{Name: "p1", File: writePatch(t, "")}},
	}

	_, err := runner.Run(cfg, runner.RunOptions{StatePath: statePath, DryRun: true})
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}

	// State file must NOT be written during a dry run.
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("state file should not exist after dry run")
	}
}

func TestRun_MissingRepo(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.yaml")
	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "missing", Path: "/nonexistent/path"}},
		Patches: []config.Patch{{Name: "p1", File: writePatch(t, "")}},
	}

	results, err := runner.Run(cfg, runner.RunOptions{StatePath: statePath})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Err == nil {
		t.Errorf("expected error result for missing repo, got %+v", results)
	}
}

// loadState is a thin helper so the test file compiles without importing state directly.
func loadState(path string) (interface{}, error) { return nil, nil }
