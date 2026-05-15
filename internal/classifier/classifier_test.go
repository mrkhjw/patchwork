package classifier_test

import (
	"testing"

	"github.com/patchwork/internal/classifier"
	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/state"
)

func sampleConfig() *config.Config {
	return &config.Config{
		Patches: []config.Patch{
			{Name: "add-index", Repo: "db", Tags: []string{"schema"}},
			{Name: "fix-xss", Repo: "web", Tags: []string{"security"}},
			{Name: "update-copy", Repo: "web", Tags: []string{}},
			{Name: "migrate-users", Repo: "db", Tags: []string{"migration"}},
		},
	}
}

func sampleState(statuses map[string]string) *state.State {
	st := state.New()
	for name, status := range statuses {
		st.Upsert(state.Entry{PatchName: name, Status: status})
	}
	return st
}

func TestClassify_CriticalTag(t *testing.T) {
	cfg := sampleConfig()
	st := sampleState(nil)
	results := classifier.Classify(cfg, st, classifier.DefaultPolicy())

	for _, r := range results {
		if r.PatchName == "fix-xss" && r.Tier != classifier.TierCritical {
			t.Errorf("expected fix-xss to be critical, got %s", r.Tier)
		}
	}
}

func TestClassify_HighTag(t *testing.T) {
	cfg := sampleConfig()
	st := sampleState(nil)
	results := classifier.Classify(cfg, st, classifier.DefaultPolicy())

	for _, r := range results {
		if r.PatchName == "add-index" && r.Tier != classifier.TierHigh {
			t.Errorf("expected add-index to be high, got %s", r.Tier)
		}
	}
}

func TestClassify_FailedStatusElevates(t *testing.T) {
	cfg := sampleConfig()
	st := sampleState(map[string]string{"update-copy": "failed"})
	results := classifier.Classify(cfg, st, classifier.DefaultPolicy())

	for _, r := range results {
		if r.PatchName == "update-copy" && r.Tier != classifier.TierHigh {
			t.Errorf("expected update-copy to be high due to failed status, got %s", r.Tier)
		}
	}
}

func TestClassify_PendingIsMedium(t *testing.T) {
	cfg := sampleConfig()
	st := sampleState(map[string]string{"update-copy": "pending"})
	results := classifier.Classify(cfg, st, classifier.DefaultPolicy())

	for _, r := range results {
		if r.PatchName == "update-copy" && r.Tier != classifier.TierMedium {
			t.Errorf("expected update-copy to be medium, got %s", r.Tier)
		}
	}
}

func TestAnyAbove_True(t *testing.T) {
	results := []classifier.Result{
		{PatchName: "a", Tier: classifier.TierCritical},
		{PatchName: "b", Tier: classifier.TierLow},
	}
	if !classifier.AnyAbove(results, classifier.TierHigh) {
		t.Error("expected AnyAbove to return true")
	}
}

func TestAnyAbove_False(t *testing.T) {
	results := []classifier.Result{
		{PatchName: "a", Tier: classifier.TierMedium},
	}
	if classifier.AnyAbove(results, classifier.TierHigh) {
		t.Error("expected AnyAbove to return false")
	}
}
