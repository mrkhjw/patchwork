package differ_test

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork/internal/differ"
	"github.com/patchwork/internal/snapshot"
)

func makeSnapshot(entries map[string]snapshot.Entry) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		CreatedAt: time.Now(),
		Entries:   entries,
	}
}

func TestCompare_NoDrift(t *testing.T) {
	base := makeSnapshot(map[string]snapshot.Entry{
		"repo-a::patch-1": {Repo: "repo-a", Patch: "patch-1", Status: "applied"},
	})
	cur := makeSnapshot(map[string]snapshot.Entry{
		"repo-a::patch-1": {Repo: "repo-a", Patch: "patch-1", Status: "applied"},
	})
	entries := differ.Compare(base, cur)
	if differ.AnyDrifted(entries) {
		t.Fatal("expected no drift")
	}
}

func TestCompare_StatusChanged(t *testing.T) {
	base := makeSnapshot(map[string]snapshot.Entry{
		"repo-a::patch-1": {Repo: "repo-a", Patch: "patch-1", Status: "applied"},
	})
	cur := makeSnapshot(map[string]snapshot.Entry{
		"repo-a::patch-1": {Repo: "repo-a", Patch: "patch-1", Status: "pending"},
	})
	entries := differ.Compare(base, cur)
	if !differ.AnyDrifted(entries) {
		t.Fatal("expected drift to be detected")
	}
	if entries[0].Old != "applied" || entries[0].New != "pending" {
		t.Errorf("unexpected old/new: %s / %s", entries[0].Old, entries[0].New)
	}
}

func TestCompare_NewEntryInCurrent(t *testing.T) {
	base := makeSnapshot(map[string]snapshot.Entry{})
	cur := makeSnapshot(map[string]snapshot.Entry{
		"repo-b::patch-2": {Repo: "repo-b", Patch: "patch-2", Status: "applied"},
	})
	entries := differ.Compare(base, cur)
	if len(entries) != 1 || !entries[0].Drifted {
		t.Fatal("expected one drifted entry for new patch")
	}
	if entries[0].Old != "(absent)" {
		t.Errorf("expected old=(absent), got %s", entries[0].Old)
	}
}

func TestCompare_MissingInCurrent(t *testing.T) {
	base := makeSnapshot(map[string]snapshot.Entry{
		"repo-c::patch-3": {Repo: "repo-c", Patch: "patch-3", Status: "applied"},
	})
	cur := makeSnapshot(map[string]snapshot.Entry{})
	entries := differ.Compare(base, cur)
	if len(entries) != 1 || !entries[0].Drifted {
		t.Fatal("expected one drifted entry for removed patch")
	}
	if entries[0].New != "(absent)" {
		t.Errorf("expected new=(absent), got %s", entries[0].New)
	}
}

func TestFormat_NoDrift(t *testing.T) {
	out := differ.Format(nil)
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestFormat_ContainsHeaders(t *testing.T) {
	entries := []differ.Entry{
		{Repo: "r", Patch: "p", Old: "applied", New: "pending", Drifted: true},
	}
	out := differ.Format(entries)
	for _, h := range []string{"REPO", "PATCH", "OLD", "NEW", "DRIFTED"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q in output", h)
		}
	}
	if !strings.Contains(out, "YES") {
		t.Error("expected YES in drifted column")
	}
}
