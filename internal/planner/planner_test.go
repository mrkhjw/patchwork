package planner_test

import (
	"testing"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/planner"
	"github.com/yourorg/patchwork/internal/state"
	"github.com/yourorg/patchwork/internal/validator"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "fix-auth", File: "fix-auth.patch", Repos: []string{"repo-a", "repo-b"}, Tags: []string{"security"}},
		{Name: "add-logging", File: "add-logging.patch", Repos: []string{"repo-a"}, Tags: []string{"ops"}},
		{Name: "bad-patch", File: "bad.patch", Repos: []string{"repo-c"}, Tags: []string{}},
	}
}

func allValid(patches []config.Patch) []validator.Result {
	results := make([]validator.Result, len(patches))
	for i, p := range patches {
		results[i] = validator.Result{PatchName: p.Name, OK: true}
	}
	return results
}

func TestPlan_AllPending(t *testing.T) {
	st := state.New()
	results := allValid(samplePatches())
	tasks := planner.Plan(samplePatches(), st, results, planner.Options{})
	// fix-auth -> 2 repos, add-logging -> 1, bad-patch -> 1
	if len(tasks) != 4 {
		t.Fatalf("expected 4 tasks, got %d", len(tasks))
	}
}

func TestPlan_SkipsInvalidPatch(t *testing.T) {
	st := state.New()
	results := []validator.Result{
		{PatchName: "fix-auth", OK: true},
		{PatchName: "add-logging", OK: false, Error: "missing file"},
		{PatchName: "bad-patch", OK: true},
	}
	tasks := planner.Plan(samplePatches(), st, results, planner.Options{})
	for _, tk := range tasks {
		if tk.Patch.Name == "add-logging" {
			t.Fatal("invalid patch should be excluded")
		}
	}
}

func TestPlan_SkipsApplied(t *testing.T) {
	st := state.New()
	st.Upsert("fix-auth", "repo-a", state.Entry{Status: "applied"})
	results := allValid(samplePatches())
	tasks := planner.Plan(samplePatches(), st, results, planner.Options{SkipApplied: true})
	for _, tk := range tasks {
		if tk.Patch.Name == "fix-auth" && tk.Repo == "repo-a" {
			t.Fatal("already-applied task should be skipped")
		}
	}
}

func TestPlan_FilterByTag(t *testing.T) {
	st := state.New()
	results := allValid(samplePatches())
	tasks := planner.Plan(samplePatches(), st, results, planner.Options{FilterTags: []string{"security"}})
	for _, tk := range tasks {
		if tk.Patch.Name != "fix-auth" {
			t.Fatalf("expected only fix-auth tasks, got %s", tk.Patch.Name)
		}
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestPlan_ReasonReflectsState(t *testing.T) {
	st := state.New()
	st.Upsert("fix-auth", "repo-b", state.Entry{Status: "failed"})
	results := allValid(samplePatches())
	tasks := planner.Plan(samplePatches(), st, results, planner.Options{FilterNames: []string{"fix-auth"}})
	for _, tk := range tasks {
		if tk.Repo == "repo-b" && tk.Reason != "failed" {
			t.Fatalf("expected reason 'failed', got %q", tk.Reason)
		}
		if tk.Repo == "repo-a" && tk.Reason != "pending" {
			t.Fatalf("expected reason 'pending', got %q", tk.Reason)
		}
	}
}
