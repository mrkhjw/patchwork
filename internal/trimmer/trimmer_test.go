package trimmer_test

import (
	"testing"
	"time"

	"github.com/patchwork/internal/state"
	"github.com/patchwork/internal/trimmer"
)

func buildState(t *testing.T, entries []state.Entry) *state.State {
	t.Helper()
	s := state.New()
	for _, e := range entries {
		s.Upsert(e)
	}
	return s
}

func TestTrim_NoPolicy_RemovesNothing(t *testing.T) {
	now := time.Now()
	s := buildState(t, []state.Entry{
		{PatchName: "alpha", Status: "applied", UpdatedAt: now.Add(-48 * time.Hour)},
		{PatchName: "beta", Status: "pending", UpdatedAt: now},
	})

	results := trimmer.Trim(s, trimmer.Policy{})
	if len(results) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(results))
	}
	if len(s.All()) != 2 {
		t.Fatalf("state should still have 2 entries")
	}
}

func TestTrim_MaxAge_RemovesOldEntries(t *testing.T) {
	now := time.Now()
	s := buildState(t, []state.Entry{
		{PatchName: "old", Status: "applied", UpdatedAt: now.Add(-72 * time.Hour)},
		{PatchName: "fresh", Status: "applied", UpdatedAt: now.Add(-1 * time.Hour)},
	})

	results := trimmer.Trim(s, trimmer.Policy{MaxAge: 48 * time.Hour})
	if len(results) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(results))
	}
	if results[0].PatchName != "old" {
		t.Errorf("expected 'old' to be removed, got %q", results[0].PatchName)
	}
	if len(s.All()) != 1 {
		t.Fatalf("state should have 1 entry remaining")
	}
}

func TestTrim_MaxEntries_KeepsNewest(t *testing.T) {
	now := time.Now()
	s := buildState(t, []state.Entry{
		{PatchName: "a", UpdatedAt: now.Add(-3 * time.Hour)},
		{PatchName: "b", UpdatedAt: now.Add(-1 * time.Hour)},
		{PatchName: "c", UpdatedAt: now.Add(-2 * time.Hour)},
	})

	results := trimmer.Trim(s, trimmer.Policy{MaxEntries: 2})
	if len(results) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(results))
	}
	if results[0].PatchName != "a" {
		t.Errorf("expected oldest 'a' removed, got %q", results[0].PatchName)
	}
}

func TestTrim_CombinedPolicy(t *testing.T) {
	now := time.Now()
	s := buildState(t, []state.Entry{
		{PatchName: "ancient", UpdatedAt: now.Add(-100 * time.Hour)},
		{PatchName: "mid", UpdatedAt: now.Add(-10 * time.Hour)},
		{PatchName: "recent", UpdatedAt: now.Add(-1 * time.Hour)},
	})

	results := trimmer.Trim(s, trimmer.Policy{MaxAge: 24 * time.Hour, MaxEntries: 1})
	// "ancient" removed by age; then "mid" removed by count.
	if len(results) != 2 {
		t.Fatalf("expected 2 removed, got %d: %+v", len(results), results)
	}
	if len(s.All()) != 1 {
		t.Fatalf("expected 1 entry remaining")
	}
}

func TestAnyRemoved_TrueAndFalse(t *testing.T) {
	if trimmer.AnyRemoved(nil) {
		t.Error("expected false for nil slice")
	}
	if !trimmer.AnyRemoved([]trimmer.Result{{PatchName: "x"}}) {
		t.Error("expected true for non-empty slice")
	}
}

func TestFormat_ContainsPatchName(t *testing.T) {
	out := trimmer.Format([]trimmer.Result{
		{PatchName: "my-patch", Reason: "exceeds max entry count"},
	})
	if !contains(out, "my-patch") {
		t.Errorf("format output missing patch name: %s", out)
	}
	if !contains(out, "removed 1") {
		t.Errorf("format output missing count: %s", out)
	}
}

func TestFormat_Empty(t *testing.T) {
	out := trimmer.Format(nil)
	if !contains(out, "nothing removed") {
		t.Errorf("expected 'nothing removed' in output: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
