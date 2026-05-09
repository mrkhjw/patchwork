package grouper_test

import (
	"testing"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/grouper"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "alpha", Repos: []string{"repo-a"}, Tags: []string{"hotfix", "db"}},
		{Name: "beta", Repos: []string{"repo-b"}, Tags: []string{"hotfix"}},
		{Name: "gamma", Repos: []string{"repo-a", "repo-c"}, Tags: []string{}},
	}
}

func TestGroup_ByRepo(t *testing.T) {
	res, err := grouper.Group(samplePatches(), grouper.ByRepo, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["repo-a"]) != 2 {
		t.Errorf("repo-a: want 2 patches, got %d", len(res.Buckets["repo-a"]))
	}
	if len(res.Buckets["repo-b"]) != 1 {
		t.Errorf("repo-b: want 1 patch, got %d", len(res.Buckets["repo-b"]))
	}
}

func TestGroup_ByTag_Untagged(t *testing.T) {
	res, err := grouper.Group(samplePatches(), grouper.ByTag, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["(untagged)"]) != 1 {
		t.Errorf("(untagged): want 1, got %d", len(res.Buckets["(untagged)"]))
	}
	if len(res.Buckets["hotfix"]) != 2 {
		t.Errorf("hotfix: want 2, got %d", len(res.Buckets["hotfix"]))
	}
}

func TestGroup_ByStatus(t *testing.T) {
	statuses := map[string]string{
		"alpha": "applied",
		"beta":  "pending",
		// gamma has no entry → "unknown"
	}
	res, err := grouper.Group(samplePatches(), grouper.ByStatus, statuses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["applied"]) != 1 {
		t.Errorf("applied: want 1, got %d", len(res.Buckets["applied"]))
	}
	if len(res.Buckets["unknown"]) != 1 {
		t.Errorf("unknown: want 1, got %d", len(res.Buckets["unknown"]))
	}
}

func TestGroup_UnknownDimension(t *testing.T) {
	_, err := grouper.Group(samplePatches(), grouper.By("invalid"), nil)
	if err == nil {
		t.Error("expected error for unknown dimension, got nil")
	}
}

func TestResult_Counts(t *testing.T) {
	res, _ := grouper.Group(samplePatches(), grouper.ByRepo, nil)
	counts := res.Counts()
	if counts["repo-a"] != 2 {
		t.Errorf("counts repo-a: want 2, got %d", counts["repo-a"])
	}
}

func TestResult_KeysAreSorted(t *testing.T) {
	res, _ := grouper.Group(samplePatches(), grouper.ByRepo, nil)
	for i := 1; i < len(res.Keys); i++ {
		if res.Keys[i] < res.Keys[i-1] {
			t.Errorf("keys not sorted at index %d: %v", i, res.Keys)
		}
	}
}
