package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/snapshot"
)

func TestSaveAndLoad_Roundtrip(t *testing.T) {
	entries := []snapshot.Entry{
		{PatchName: "fix-lint", Repo: "repo-a", Applied: true},
		{PatchName: "add-ci", Repo: "repo-b", Applied: false},
	}
	s := snapshot.New(entries)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(tmp, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].PatchName != "fix-lint" {
		t.Errorf("unexpected patch name: %s", loaded.Entries[0].PatchName)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_StampsTime(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New([]snapshot.Entry{{PatchName: "p", Repo: "r", Applied: true}})
	after := time.Now().UTC()

	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v out of expected range", s.CapturedAt)
	}
	for _, e := range s.Entries {
		if e.CapturedAt.IsZero() {
			t.Errorf("entry CapturedAt not set for %s", e.PatchName)
		}
	}
}

func TestDiff_DetectsDrift(t *testing.T) {
	baseline := snapshot.New([]snapshot.Entry{
		{PatchName: "fix-lint", Repo: "repo-a", Applied: true},
		{PatchName: "add-ci", Repo: "repo-b", Applied: false},
	})
	current := snapshot.New([]snapshot.Entry{
		{PatchName: "fix-lint", Repo: "repo-a", Applied: true},  // unchanged
		{PatchName: "add-ci", Repo: "repo-b", Applied: true},   // drifted
		{PatchName: "new-patch", Repo: "repo-c", Applied: true}, // new
	})

	drifted := snapshot.Diff(baseline, current)
	if len(drifted) != 2 {
		t.Errorf("expected 2 drifted entries, got %d", len(drifted))
	}
}

func TestDiff_NoDrift(t *testing.T) {
	entries := []snapshot.Entry{
		{PatchName: "fix-lint", Repo: "repo-a", Applied: true},
	}
	baseline := snapshot.New(entries)
	current := snapshot.New([]snapshot.Entry{
		{PatchName: "fix-lint", Repo: "repo-a", Applied: true},
	})

	drifted := snapshot.Diff(baseline, current)
	if len(drifted) != 0 {
		t.Errorf("expected no drift, got %d entries", len(drifted))
	}
}

func TestSave_CreatesFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.json")
	s := snapshot.New(nil)
	if err := snapshot.Save(tmp, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
