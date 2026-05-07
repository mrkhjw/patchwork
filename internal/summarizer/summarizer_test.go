package summarizer_test

import (
	"strings"
	"testing"
	"time"

	"patchwork/internal/state"
	"patchwork/internal/summarizer"
)

func sampleEntries() []state.Entry {
	now := time.Now().UTC()
	return []state.Entry{
		{PatchName: "patch-a", Repo: "repo1", Status: "applied", AppliedAt: now},
		{PatchName: "patch-b", Repo: "repo1", Status: "skipped", AppliedAt: now},
		{PatchName: "patch-c", Repo: "repo2", Status: "failed", AppliedAt: now},
		{PatchName: "patch-d", Repo: "repo2", Status: "applied", AppliedAt: now},
	}
}

func TestBuild_Counts(t *testing.T) {
	s := summarizer.Build(sampleEntries())

	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Applied != 2 {
		t.Errorf("expected Applied=2, got %d", s.Applied)
	}
	if s.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", s.Skipped)
	}
	if s.Failed != 1 {
		t.Errorf("expected Failed=1, got %d", s.Failed)
	}
}

func TestBuild_EmptyEntries(t *testing.T) {
	s := summarizer.Build([]state.Entry{})

	if s.Total != 0 || s.Applied != 0 || s.Skipped != 0 || s.Failed != 0 {
		t.Error("expected all zero counts for empty entries")
	}
}

func TestBuild_PatchesPopulated(t *testing.T) {
	s := summarizer.Build(sampleEntries())

	if len(s.Patches) != 4 {
		t.Fatalf("expected 4 patch summaries, got %d", len(s.Patches))
	}
	if s.Patches[0].PatchName != "patch-a" {
		t.Errorf("unexpected patch name: %s", s.Patches[0].PatchName)
	}
}

func TestFormat_ContainsKeyFields(t *testing.T) {
	s := summarizer.Build(sampleEntries())
	out := summarizer.Format(s)

	for _, keyword := range []string{"Total", "Applied", "Skipped", "Failed", "Run at"} {
		if !strings.Contains(out, keyword) {
			t.Errorf("expected Format output to contain %q", keyword)
		}
	}
}

func TestAnyFailed_True(t *testing.T) {
	s := summarizer.Build(sampleEntries())
	if !summarizer.AnyFailed(s) {
		t.Error("expected AnyFailed=true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	entries := []state.Entry{
		{PatchName: "patch-a", Repo: "repo1", Status: "applied", AppliedAt: time.Now()},
	}
	s := summarizer.Build(entries)
	if summarizer.AnyFailed(s) {
		t.Error("expected AnyFailed=false")
	}
}
