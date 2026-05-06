package filter_test

import (
	"testing"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/filter"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "fix-auth", Repo: "api", Tags: []string{"security", "hotfix"}},
		{Name: "add-logging", Repo: "api", Tags: []string{"observability"}},
		{Name: "bump-deps", Repo: "frontend", Tags: []string{"maintenance"}},
		{Name: "fix-cors", Repo: "frontend", Tags: []string{"security"}},
	}
}

func TestApply_NoFilter(t *testing.T) {
	result := filter.Apply(samplePatches(), filter.Options{})
	if len(result) != 4 {
		t.Fatalf("expected 4 patches, got %d", len(result))
	}
}

func TestApply_FilterByName(t *testing.T) {
	opts := filter.Options{Names: []string{"fix-auth", "bump-deps"}}
	result := filter.Apply(samplePatches(), opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Name != "fix-auth" || result[1].Name != "bump-deps" {
		t.Errorf("unexpected patch names: %v", result)
	}
}

func TestApply_FilterByTag(t *testing.T) {
	opts := filter.Options{Tags: []string{"security"}}
	result := filter.Apply(samplePatches(), opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 security patches, got %d", len(result))
	}
}

func TestApply_FilterByRepo(t *testing.T) {
	opts := filter.Options{Repos: []string{"frontend"}}
	result := filter.Apply(samplePatches(), opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 frontend patches, got %d", len(result))
	}
}

func TestApply_CombinedFilter(t *testing.T) {
	opts := filter.Options{
		Repos: []string{"api"},
		Tags:  []string{"security"},
	}
	result := filter.Apply(samplePatches(), opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Name != "fix-auth" {
		t.Errorf("expected fix-auth, got %s", result[0].Name)
	}
}

func TestApply_NoMatch(t *testing.T) {
	opts := filter.Options{Names: []string{"nonexistent"}}
	result := filter.Apply(samplePatches(), opts)
	if len(result) != 0 {
		t.Fatalf("expected 0 results, got %d", len(result))
	}
}
