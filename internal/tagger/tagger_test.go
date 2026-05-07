package tagger_test

import (
	"testing"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/tagger"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "alpha", Tags: []string{"security", "hotfix"}},
		{Name: "beta", Tags: []string{"security"}},
		{Name: "gamma", Tags: []string{"hotfix", "experimental"}},
		{Name: "delta", Tags: []string{}},
	}
}

func TestBuild_IndexContainsAllTags(t *testing.T) {
	idx := tagger.Build(samplePatches())
	got := idx.AllTags()
	want := []string{"experimental", "hotfix", "security"}
	if len(got) != len(want) {
		t.Fatalf("AllTags() len = %d, want %d", len(got), len(want))
	}
	for i, tag := range want {
		if got[i] != tag {
			t.Errorf("AllTags()[%d] = %q, want %q", i, got[i], tag)
		}
	}
}

func TestBuild_PatchesForTag(t *testing.T) {
	idx := tagger.Build(samplePatches())

	got := idx.PatchesForTag("security")
	if len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Errorf("PatchesForTag(security) = %v, want [alpha beta]", got)
	}

	got = idx.PatchesForTag("hotfix")
	if len(got) != 2 || got[0] != "alpha" || got[1] != "gamma" {
		t.Errorf("PatchesForTag(hotfix) = %v, want [alpha gamma]", got)
	}
}

func TestBuild_UnknownTag(t *testing.T) {
	idx := tagger.Build(samplePatches())
	if got := idx.PatchesForTag("nonexistent"); got != nil {
		t.Errorf("expected nil for unknown tag, got %v", got)
	}
}

func TestBuild_PatchWithNoTags(t *testing.T) {
	idx := tagger.Build(samplePatches())
	for _, patches := range idx {
		for _, name := range patches {
			if name == "delta" {
				t.Error("patch 'delta' with no tags should not appear in index")
			}
		}
	}
}

func TestTagsForPatch_Found(t *testing.T) {
	got := tagger.TagsForPatch(samplePatches(), "alpha")
	if len(got) != 2 || got[0] != "hotfix" || got[1] != "security" {
		t.Errorf("TagsForPatch(alpha) = %v, want [hotfix security]", got)
	}
}

func TestTagsForPatch_NoTags(t *testing.T) {
	got := tagger.TagsForPatch(samplePatches(), "delta")
	if got != nil {
		t.Errorf("expected nil for patch with no tags, got %v", got)
	}
}

func TestTagsForPatch_NotFound(t *testing.T) {
	got := tagger.TagsForPatch(samplePatches(), "unknown")
	if got != nil {
		t.Errorf("expected nil for unknown patch, got %v", got)
	}
}
