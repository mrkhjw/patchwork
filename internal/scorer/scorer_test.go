package scorer_test

import (
	"testing"
	"time"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/scorer"
	"github.com/patchwork/internal/state"
)

func sampleConfig() *config.Config {
	return &config.Config{
		Patches: []config.Patch{
			{Name: "alpha", Tags: []string{"critical"}, File: "a.patch", Repos: []string{"r1"}},
			{Name: "beta", Tags: []string{"hotfix"}, File: "b.patch", Repos: []string{"r1"}},
			{Name: "gamma", Tags: []string{}, File: "c.patch", Repos: []string{"r1"}},
		},
	}
}

func emptyState(t *testing.T) *state.State {
	t.Helper()
	st, err := state.New()
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}
	return st
}

func TestScore_TagWeights(t *testing.T) {
	cfg := sampleConfig()
	st := emptyState(t)
	pol := scorer.DefaultPolicy()

	results := scorer.Score(cfg, st, pol)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	byName := map[string]float64{}
	for _, r := range results {
		byName[r.PatchName] = r.Score
	}

	if byName["alpha"] != 10.0 {
		t.Errorf("alpha score: want 10.0, got %f", byName["alpha"])
	}
	if byName["beta"] != 5.0 {
		t.Errorf("beta score: want 5.0, got %f", byName["beta"])
	}
	if byName["gamma"] != 0.0 {
		t.Errorf("gamma score: want 0.0, got %f", byName["gamma"])
	}
}

func TestScore_RetryPenalty(t *testing.T) {
	cfg := sampleConfig()
	st := emptyState(t)
	st.Upsert(state.Entry{PatchName: "alpha", Attempts: 6, AppliedAt: time.Time{}})

	pol := scorer.DefaultPolicy() // RetryPenalty = 2.0
	results := scorer.Score(cfg, st, pol)

	byName := map[string]float64{}
	for _, r := range results {
		byName[r.PatchName] = r.Score
	}

	// 10 (critical) - 2*6 = -2
	if byName["alpha"] != -2.0 {
		t.Errorf("alpha score: want -2.0, got %f", byName["alpha"])
	}
}

func TestAnyNegative_True(t *testing.T) {
	results := []scorer.Result{
		{PatchName: "a", Score: 3.0},
		{PatchName: "b", Score: -1.0},
	}
	if !scorer.AnyNegative(results) {
		t.Error("expected AnyNegative to return true")
	}
}

func TestAnyNegative_False(t *testing.T) {
	results := []scorer.Result{
		{PatchName: "a", Score: 0.0},
		{PatchName: "b", Score: 5.0},
	}
	if scorer.AnyNegative(results) {
		t.Error("expected AnyNegative to return false")
	}
}

func TestScore_EmptyConfig(t *testing.T) {
	cfg := &config.Config{}
	st := emptyState(t)
	results := scorer.Score(cfg, st, scorer.DefaultPolicy())
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
