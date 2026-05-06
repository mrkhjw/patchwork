package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork/internal/audit"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.json")
}

func TestAppendAndFilter(t *testing.T) {
	l := &audit.Log{}
	l.Append("fix-typo", "repo-a", audit.EventApplied, "")
	l.Append("add-ci", "repo-b", audit.EventSkipped, "already applied")
	l.Append("fix-typo", "repo-c", audit.EventFailed, "patch rejected")

	if len(l.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(l.Events))
	}

	applied := l.Filter(audit.EventApplied)
	if len(applied) != 1 || applied[0].Patch != "fix-typo" {
		t.Errorf("unexpected applied filter result: %+v", applied)
	}

	all := l.Filter("")
	if len(all) != 3 {
		t.Errorf("expected all 3 events, got %d", len(all))
	}
}

func TestSaveAndLoad_Roundtrip(t *testing.T) {
	path := tempPath(t)

	l := &audit.Log{}
	l.Append("patch-x", "repo-1", audit.EventRolledBack, "user requested")
	l.Append("patch-y", "repo-2", audit.EventApplied, "")

	if err := l.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Events) != 2 {
		t.Fatalf("expected 2 events after reload, got %d", len(loaded.Events))
	}
	if loaded.Events[0].Patch != "patch-x" || loaded.Events[0].Kind != audit.EventRolledBack {
		t.Errorf("unexpected first event: %+v", loaded.Events[0])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	l, err := audit.Load(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(l.Events) != 0 {
		t.Errorf("expected empty log, got %d events", len(l.Events))
	}
}

func TestEvent_TimestampSet(t *testing.T) {
	l := &audit.Log{}
	l.Append("p", "r", audit.EventApplied, "")
	if l.Events[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	path := tempPath(t)
	l := &audit.Log{}
	l.Append("p", "r", audit.EventSkipped, "")
	if err := l.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
