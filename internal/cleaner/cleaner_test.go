package cleaner_test

import (
	"os"
	"path/filepath"
	"testing"

	"patchwork/internal/cleaner"
	"patchwork/internal/config"
	"patchwork/internal/state"
)

func buildConfig(patches []config.Patch) *config.Config {
	return &config.Config{Patches: patches}
}

func buildState(entries []state.Entry) *state.State {
	st := state.New()
	for _, e := range entries {
		st.Upsert(e)
	}
	return st
}

func TestClean_RemovesAppliedWhenFileGone(t *testing.T) {
	cfg := buildConfig([]config.Patch{
		{Name: "fix-a", File: "/nonexistent/fix-a.patch", Repos: []string{"repo1"}},
	})
	st := buildState([]state.Entry{
		{PatchName: "fix-a", Repo: "repo1", Status: "applied"},
	})

	results := cleaner.Clean(cfg, st)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Removed {
		t.Error("expected entry to be removed")
	}
}

func TestClean_KeepsAppliedWhenFileExists(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "fix-b.patch")
	_ = os.WriteFile(f, []byte("diff"), 0o644)

	cfg := buildConfig([]config.Patch{
		{Name: "fix-b", File: f, Repos: []string{"repo1"}},
	})
	st := buildState([]state.Entry{
		{PatchName: "fix-b", Repo: "repo1", Status: "applied"},
	})

	results := cleaner.Clean(cfg, st)

	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestClean_SkipsPendingEvenIfMissingFromConfig(t *testing.T) {
	cfg := buildConfig([]config.Patch{})
	st := buildState([]state.Entry{
		{PatchName: "orphan", Repo: "repo1", Status: "pending"},
	})

	results := cleaner.Clean(cfg, st)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Removed {
		t.Error("pending entry should not be removed")
	}
}

func TestClean_RemovesAppliedMissingFromConfig(t *testing.T) {
	cfg := buildConfig([]config.Patch{})
	st := buildState([]state.Entry{
		{PatchName: "old-patch", Repo: "repo2", Status: "applied"},
	})

	results := cleaner.Clean(cfg, st)

	if !cleaner.AnyRemoved(results) {
		t.Error("expected AnyRemoved to be true")
	}
}

func TestAnyRemoved_False(t *testing.T) {
	results := []cleaner.Result{
		{PatchName: "x", Removed: false},
	}
	if cleaner.AnyRemoved(results) {
		t.Error("expected AnyRemoved to be false")
	}
}
