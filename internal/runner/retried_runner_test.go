package runner

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/retrier"
)

func TestRunWithRetry_SuccessFirstAttempt(t *testing.T) {
	repoDir := initRepo(t)
	patchFile := writePatch(t, repoDir)

	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r", Path: repoDir}},
		Patches: []config.Patch{{Name: "p", Path: patchFile, Repos: []string{"r"}}},
	}
	opts := RetriedOptions{
		Options: Options{DryRun: false, StateDir: t.TempDir()},
		Retry:   retrier.Policy{MaxAttempts: 3, Delay: 0, Backoff: 1.0},
	}

	results := RunWithRetry(cfg, opts)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
}

func TestRunWithRetry_MissingPatchFile(t *testing.T) {
	repoDir := initRepo(t)

	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r", Path: repoDir}},
		Patches: []config.Patch{{Name: "p", Path: "/nonexistent/patch.diff", Repos: []string{"r"}}},
	}
	opts := RetriedOptions{
		Options: Options{DryRun: false, StateDir: t.TempDir()},
		Retry:   retrier.Policy{MaxAttempts: 2, Delay: 0, Backoff: 1.0},
	}

	results := RunWithRetry(cfg, opts)
	if len(results) == 0 {
		t.Fatal("expected a result even on failure")
	}
}

func TestRunWithRetry_DryRunNoSideEffects(t *testing.T) {
	repoDir := initRepo(t)
	patchFile := writePatch(t, repoDir)
	stateDir := t.TempDir()

	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r", Path: repoDir}},
		Patches: []config.Patch{{Name: "p", Path: patchFile, Repos: []string{"r"}}},
	}
	opts := RetriedOptions{
		Options: Options{DryRun: true, StateDir: stateDir},
		Retry:   retrier.Policy{MaxAttempts: 1, Delay: 0, Backoff: 1.0},
	}

	RunWithRetry(cfg, opts)

	// State file must not be written in dry-run mode.
	entries, _ := os.ReadDir(stateDir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			t.Fatalf("state file written during dry run: %s", e.Name())
		}
	}
}

func TestRunWithRetry_ZeroDelayPolicy(t *testing.T) {
	repoDir := initRepo(t)
	patchFile := writePatch(t, repoDir)

	start := time.Now()
	cfg := &config.Config{
		Repos:   []config.Repo{{Name: "r", Path: repoDir}},
		Patches: []config.Patch{{Name: "p", Path: patchFile, Repos: []string{"r"}}},
	}
	opts := RetriedOptions{
		Options: Options{DryRun: false, StateDir: t.TempDir()},
		Retry:   retrier.Policy{MaxAttempts: 3, Delay: 0, Backoff: 1.0},
	}
	RunWithRetry(cfg, opts)
	if time.Since(start) > 2*time.Second {
		t.Fatal("zero-delay retry took too long")
	}
}
