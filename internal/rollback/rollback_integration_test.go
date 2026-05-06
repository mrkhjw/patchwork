package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/rollback"
	"github.com/patchwork/internal/state"
)

// TestRoundtrip_MultiplePatches verifies that independent backups for several
// patches coexist without interfering with each other.
func TestRoundtrip_MultiplePatches(t *testing.T) {
	repo := setupDir(t)

	files := map[string]string{
		"patch-x": "content-x",
		"patch-y": "content-y",
		"patch-z": "content-z",
	}

	targets := make(map[string]string)
	for name, content := range files {
		path := filepath.Join(repo, name+".txt")
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		targets[name] = path
		if err := rollback.Backup(repo, name, path); err != nil {
			t.Fatalf("Backup %s: %v", name, err)
		}
	}

	// Overwrite all files to simulate patches.
	for name, path := range targets {
		if err := os.WriteFile(path, []byte("patched-"+name), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	st := state.New()
	for name := range files {
		st.Upsert(state.Entry{Name: name, Status: state.StatusApplied})
	}

	// Restore each and verify original content.
	for name, path := range targets {
		if err := rollback.Restore(repo, name, path, st); err != nil {
			t.Fatalf("Restore %s: %v", name, err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != files[name] {
			t.Errorf("%s: expected %q, got %q", name, files[name], string(got))
		}
		entry, _ := st.Get(name)
		if entry.Status != state.StatusPending {
			t.Errorf("%s: expected pending status after restore", name)
		}
		if rollback.HasBackup(repo, name) {
			t.Errorf("%s: backup should be removed after restore", name)
		}
	}
}
