package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/rollback"
	"github.com/patchwork/internal/state"
)

func setupDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rollback-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestBackupAndRestore_ExistingFile(t *testing.T) {
	repo := setupDir(t)
	target := filepath.Join(repo, "hello.txt")
	original := []byte("original content\n")
	if err := os.WriteFile(target, original, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := rollback.Backup(repo, "patch-a", target); err != nil {
		t.Fatalf("Backup: %v", err)
	}

	// Simulate patch modifying the file.
	if err := os.WriteFile(target, []byte("patched content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	st := state.New()
	st.Upsert(state.Entry{Name: "patch-a", Status: state.StatusApplied})

	if err := rollback.Restore(repo, "patch-a", target, st); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	got, _ := os.ReadFile(target)
	if string(got) != string(original) {
		t.Errorf("expected original content, got %q", string(got))
	}

	entry, _ := st.Get("patch-a")
	if entry.Status != state.StatusPending {
		t.Errorf("expected status pending, got %q", entry.Status)
	}
}

func TestBackupAndRestore_NonExistentFile(t *testing.T) {
	repo := setupDir(t)
	target := filepath.Join(repo, "new.txt")

	if err := rollback.Backup(repo, "patch-b", target); err != nil {
		t.Fatalf("Backup: %v", err)
	}

	if err := os.WriteFile(target, []byte("created by patch\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	st := state.New()
	if err := rollback.Restore(repo, "patch-b", target, st); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Error("expected file to be removed after restore")
	}
}

func TestHasBackup(t *testing.T) {
	repo := setupDir(t)
	target := filepath.Join(repo, "f.txt")
	_ = os.WriteFile(target, []byte("x"), 0o644)

	if rollback.HasBackup(repo, "patch-c") {
		t.Error("expected no backup before Backup()")
	}
	_ = rollback.Backup(repo, "patch-c", target)
	if !rollback.HasBackup(repo, "patch-c") {
		t.Error("expected backup after Backup()")
	}
}

func TestRestore_MissingBackup(t *testing.T) {
	repo := setupDir(t)
	target := filepath.Join(repo, "missing.txt")
	st := state.New()

	err := rollback.Restore(repo, "no-such-patch", target, st)
	if err == nil {
		t.Error("expected error when backup does not exist")
	}
}
