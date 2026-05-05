package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestLoad_MissingFile(t *testing.T) {
	s, err := Load("/nonexistent/path/state.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected empty state, got %d entries", len(s.Entries))
	}
}

func TestSaveAndLoad_Roundtrip(t *testing.T) {
	path := tempPath(t)
	s := &State{
		Entries: []PatchStatus{
			{PatchName: "fix-timeout", Repo: "api", Applied: true, AppliedAt: time.Now().UTC().Truncate(time.Second), Checksum: "abc123"},
		},
	}
	if err := s.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].PatchName != "fix-timeout" {
		t.Errorf("unexpected patch name: %s", loaded.Entries[0].PatchName)
	}
}

func TestUpsert_AddsAndUpdates(t *testing.T) {
	s := &State{}
	ps := PatchStatus{PatchName: "p1", Repo: "r1", Applied: false, Checksum: "aaa"}
	s.Upsert(ps)
	if len(s.Entries) != 1 {
		t.Fatalf("expected 1 entry after insert")
	}
	ps.Applied = true
	ps.Checksum = "bbb"
	s.Upsert(ps)
	if len(s.Entries) != 1 {
		t.Fatalf("expected still 1 entry after update")
	}
	if !s.Entries[0].Applied || s.Entries[0].Checksum != "bbb" {
		t.Errorf("entry was not updated correctly")
	}
}

func TestGet_FoundAndNotFound(t *testing.T) {
	s := &State{
		Entries: []PatchStatus{
			{PatchName: "p1", Repo: "r1", Applied: true, Checksum: "xyz"},
		},
	}
	entry, ok := s.Get("p1", "r1")
	if !ok {
		t.Fatal("expected to find entry")
	}
	if entry.Checksum != "xyz" {
		t.Errorf("unexpected checksum: %s", entry.Checksum)
	}
	_, ok = s.Get("p1", "missing-repo")
	if ok {
		t.Error("expected not found for unknown repo")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
