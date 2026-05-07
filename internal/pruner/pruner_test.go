package pruner_test

import (
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/pruner"
	"github.com/patchwork/internal/state"
)

func buildState(t *testing.T, names ...string) *state.State {
	t.Helper()
	st := state.New()
	for _, n := range names {
		st.Upsert(state.Entry{PatchName: n, Status: "applied"})
	}
	return st
}

func buildConfig(names ...string) *config.Config {
	cfg := &config.Config{}
	for _, n := range names {
		cfg.Patches = append(cfg.Patches, config.Patch{Name: n})
	}
	return cfg
}

func TestPrune_RemovesStaleEntries(t *testing.T) {
	st := buildState(t, "alpha", "beta", "gamma")
	cfg := buildConfig("alpha") // beta and gamma are stale

	results, err := pruner.Prune(st, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 removals, got %d", len(results))
	}
	for _, r := range results {
		if !r.Removed {
			t.Errorf("expected Removed=true for %q", r.PatchName)
		}
	}
}

func TestPrune_KeepsValidEntries(t *testing.T) {
	st := buildState(t, "alpha", "beta")
	cfg := buildConfig("alpha", "beta")

	results, err := pruner.Prune(st, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 removals, got %d", len(results))
	}
}

func TestPrune_EmptyState(t *testing.T) {
	st := buildState(t)
	cfg := buildConfig("alpha")

	results, err := pruner.Prune(st, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 removals, got %d", len(results))
	}
}

func TestAnyRemoved_True(t *testing.T) {
	results := []pruner.Result{{PatchName: "x", Removed: true}}
	if !pruner.AnyRemoved(results) {
		t.Error("expected AnyRemoved to return true")
	}
}

func TestAnyRemoved_False(t *testing.T) {
	results := []pruner.Result{{PatchName: "x", Removed: false}}
	if pruner.AnyRemoved(results) {
		t.Error("expected AnyRemoved to return false")
	}
}
