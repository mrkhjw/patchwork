package archiver_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/archiver"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestArchive_CreatesZip(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	a := writeTempFile(t, src, "audit.log", "event1\nevent2")
	b := writeTempFile(t, src, "diff.patch", "--- a\n+++ b")

	results, err := archiver.Archive(dest, "run", []string{a, b})
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("result error for %s: %v", r.Source, r.Err)
		}
	}
}

func TestArchive_ZipContainsFiles(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeTempFile(t, src, "notes.txt", "hello")
	writeTempFile(t, src, "state.json", `{}`)

	results, err := archiver.Archive(dest, "snap", []string{
		filepath.Join(src, "notes.txt"),
		filepath.Join(src, "state.json"),
	})
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}

	zr, err := zip.OpenReader(results[0].Archive)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	names := map[string]bool{}
	for _, f := range zr.File {
		names[f.Name] = true
	}
	for _, want := range []string{"notes.txt", "state.json"} {
		if !names[want] {
			t.Errorf("zip missing entry %q", want)
		}
	}
}

func TestArchive_MissingSourceFile(t *testing.T) {
	dest := t.TempDir()
	results, err := archiver.Archive(dest, "bad", []string{"/no/such/file.log"})
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if len(results) != 1 || results[0].Err == nil {
		t.Error("expected a per-file error for missing source")
	}
	if !archiver.AnyFailed(results) {
		t.Error("AnyFailed should be true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := []archiver.Result{
		{Source: "a", Err: nil},
		{Source: "b", Err: nil},
	}
	if archiver.AnyFailed(results) {
		t.Error("AnyFailed should be false when no errors")
	}
}
