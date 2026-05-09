package matcher_test

import (
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/matcher"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "fix-auth", Repo: "/repos/api", Tags: []string{"security", "backend"}},
		{Name: "fix-logging", Repo: "/repos/api", Tags: []string{"observability"}},
		{Name: "update-deps", Repo: "/repos/frontend", Tags: []string{"maintenance"}},
		{Name: "hotfix-db", Repo: "/repos/db", Tags: []string{"backend", "urgent"}},
	}
}

func TestMatch_EmptyQuery_ReturnsAll(t *testing.T) {
	results := matcher.Match(samplePatches(), "")
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}

func TestMatch_ExactName_ScoresHighest(t *testing.T) {
	results := matcher.Match(samplePatches(), "fix-auth")
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].Patch.Name != "fix-auth" {
		t.Errorf("expected fix-auth first, got %s", results[0].Patch.Name)
	}
	if results[0].Score != 100 {
		t.Errorf("expected score 100, got %d", results[0].Score)
	}
}

func TestMatch_PartialName_ReturnsMatches(t *testing.T) {
	results := matcher.Match(samplePatches(), "fix")
	if len(results) != 2 {
		t.Fatalf("expected 2 results for 'fix', got %d", len(results))
	}
}

func TestMatch_RepoMatch(t *testing.T) {
	results := matcher.Match(samplePatches(), "frontend")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'frontend', got %d", len(results))
	}
	if results[0].Patch.Name != "update-deps" {
		t.Errorf("unexpected patch: %s", results[0].Patch.Name)
	}
}

func TestMatch_TagMatch(t *testing.T) {
	results := matcher.Match(samplePatches(), "urgent")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for tag 'urgent', got %d", len(results))
	}
	if results[0].Patch.Name != "hotfix-db" {
		t.Errorf("unexpected patch: %s", results[0].Patch.Name)
	}
}

func TestMatch_NoMatch_ReturnsEmpty(t *testing.T) {
	results := matcher.Match(samplePatches(), "zzznomatch")
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestAnyMatch_True(t *testing.T) {
	if !matcher.AnyMatch(samplePatches(), "auth") {
		t.Error("expected AnyMatch to return true")
	}
}

func TestAnyMatch_False(t *testing.T) {
	if matcher.AnyMatch(samplePatches(), "zzznomatch") {
		t.Error("expected AnyMatch to return false")
	}
}

func TestMatch_SortedByScoreDescending(t *testing.T) {
	results := matcher.Match(samplePatches(), "fix")
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted: index %d score %d > index %d score %d",
				i, results[i].Score, i-1, results[i-1].Score)
		}
	}
}
