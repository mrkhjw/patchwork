package pinner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/pinner"
)

func writeTempPatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.patch")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestPinFile_RecordsHash(t *testing.T) {
	path := writeTempPatch(t, "--- a/foo\n+++ b/foo\n@@ -1 +1 @@\n-old\n+new\n")
	store := make(pinner.Store)
	if err := pinner.PinFile(store, "mypatch", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := store["mypatch"]; !ok {
		t.Fatal("expected pin entry for 'mypatch'")
	}
	if store["mypatch"].Hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestPinFile_MissingFile(t *testing.T) {
	store := make(pinner.Store)
	err := pinner.PinFile(store, "ghost", "/nonexistent/patch.diff")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestVerify_NoDrift(t *testing.T) {
	path := writeTempPatch(t, "stable content")
	store := make(pinner.Store)
	_ = pinner.PinFile(store, "stable", path)

	results := pinner.Verify(store)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Error("expected no drift")
	}
}

func TestVerify_DetectsDrift(t *testing.T) {
	path := writeTempPatch(t, "original content")
	store := make(pinner.Store)
	_ = pinner.PinFile(store, "drifted", path)

	// Mutate the file after pinning.
	if err := os.WriteFile(path, []byte("modified content"), 0644); err != nil {
		t.Fatal(err)
	}

	results := pinner.Verify(store)
	if !results[0].Drifted {
		t.Error("expected drift to be detected")
	}
	if pinner.AnyDrifted(results) == false {
		t.Error("AnyDrifted should return true")
	}
}

func TestVerify_MissingFile(t *testing.T) {
	store := make(pinner.Store)
	store["gone"] = pinner.Pin{Name: "gone", Path: "/no/such/file", Hash: "abc123"}
	results := pinner.Verify(store)
	if !results[0].Drifted {
		t.Error("expected drift when file is missing")
	}
	if results[0].Error == "" {
		t.Error("expected error message")
	}
}

func TestSaveAndLoad_Roundtrip(t *testing.T) {
	path := writeTempPatch(t, "roundtrip content")
	store := make(pinner.Store)
	_ = pinner.PinFile(store, "rt", path)

	storePath := filepath.Join(t.TempDir(), "pins.json")
	if err := pinner.SaveStore(store, storePath); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := pinner.LoadStore(storePath)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["rt"].Hash != store["rt"].Hash {
		t.Errorf("hash mismatch: got %q want %q", loaded["rt"].Hash, store["rt"].Hash)
	}
}

func TestLoadStore_MissingFile(t *testing.T) {
	store, err := pinner.LoadStore("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store) != 0 {
		t.Error("expected empty store for missing file")
	}
}
