package labeler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork/internal/labeler"
)

func samplePatches() []labeler.PatchMeta {
	return []labeler.PatchMeta{
		{Name: "fix-auth", Repo: "github.com/org/auth-service", Status: "applied", AppliedAt: time.Now().AddDate(0, 0, -10)},
		{Name: "bump-deps", Repo: "github.com/org/frontend", Status: "pending", AppliedAt: time.Time{}},
		{Name: "hotfix-db", Repo: "github.com/org/db-layer", Status: "applied", AppliedAt: time.Now().AddDate(0, 0, -1)},
	}
}

func TestLabel_NoRules(t *testing.T) {
	results := labeler.Label(samplePatches(), nil)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if len(r.Labels) != 0 {
			t.Errorf("%s: expected no labels, got %v", r.Patch.Name, r.Labels)
		}
	}
}

func TestLabel_ByRepoContains(t *testing.T) {
	rules := []labeler.Rule{{Label: "frontend", RepoContains: "frontend"}}
	results := labeler.Label(samplePatches(), rules)
	for _, r := range results {
		if r.Patch.Name == "bump-deps" {
			if len(r.Labels) != 1 || r.Labels[0] != "frontend" {
				t.Errorf("expected [frontend], got %v", r.Labels)
			}
		} else if len(r.Labels) != 0 {
			t.Errorf("%s: unexpected labels %v", r.Patch.Name, r.Labels)
		}
	}
}

func TestLabel_ByStatus(t *testing.T) {
	rules := []labeler.Rule{{Label: "done", StatusEquals: "applied"}}
	results := labeler.Label(samplePatches(), rules)
	for _, r := range results {
		if r.Patch.Status == "applied" {
			if len(r.Labels) != 1 || r.Labels[0] != "done" {
				t.Errorf("%s: expected [done], got %v", r.Patch.Name, r.Labels)
			}
		}
	}
}

func TestLabel_ByOlderThanDays(t *testing.T) {
	rules := []labeler.Rule{{Label: "stale", OlderThanDays: 5}}
	results := labeler.Label(samplePatches(), rules)
	for _, r := range results {
		if r.Patch.Name == "fix-auth" {
			if len(r.Labels) != 1 || r.Labels[0] != "stale" {
				t.Errorf("expected [stale], got %v", r.Labels)
			}
		}
		if r.Patch.Name == "hotfix-db" && len(r.Labels) != 0 {
			t.Errorf("hotfix-db should not be stale, got %v", r.Labels)
		}
	}
}

func TestLabel_MultipleRulesNoDuplicates(t *testing.T) {
	rules := []labeler.Rule{
		{Label: "done", StatusEquals: "applied"},
		{Label: "done", RepoContains: "auth"},
	}
	results := labeler.Label(samplePatches(), rules)
	for _, r := range results {
		if r.Patch.Name == "fix-auth" {
			count := 0
			for _, l := range r.Labels {
				if l == "done" {
					count++
				}
			}
			if count != 1 {
				t.Errorf("expected label 'done' exactly once, got %d", count)
			}
		}
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := labeler.Label(samplePatches(), []labeler.Rule{{Label: "ok", StatusEquals: "applied"}})
	if labeler.AnyFailed(results) {
		t.Error("expected AnyFailed to be false")
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := []labeler.Result{
		{Patch: labeler.PatchMeta{Name: "bad"}, Labels: []string{"error"}},
	}
	if !labeler.AnyFailed(results) {
		t.Error("expected AnyFailed to be true")
	}
}

func TestFormat_WithLabels(t *testing.T) {
	r := labeler.Result{Patch: labeler.PatchMeta{Name: "fix-auth"}, Labels: []string{"done", "stale"}}
	out := labeler.Format(r)
	if !strings.Contains(out, "fix-auth") || !strings.Contains(out, "done") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestFormat_NoLabels(t *testing.T) {
	r := labeler.Result{Patch: labeler.PatchMeta{Name: "bump-deps"}, Labels: nil}
	out := labeler.Format(r)
	if !strings.Contains(out, "no labels") {
		t.Errorf("expected 'no labels' in output, got: %s", out)
	}
}
