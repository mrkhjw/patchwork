package labeler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/labeler"
)

func writeTempRules(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempRules: %v", err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	path := writeTempRules(t, `rules:
  - label: stale
    older_than_days: 30
  - label: infra
    repo_contains: infra
`)
	rules, err := labeler.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Label != "stale" {
		t.Errorf("expected label 'stale', got %q", rules[0].Label)
	}
	if rules[1].RepoContains != "infra" {
		t.Errorf("expected RepoContains 'infra', got %q", rules[1].RepoContains)
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := labeler.LoadRules("/nonexistent/rules.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRules_EmptyLabel(t *testing.T) {
	path := writeTempRules(t, `rules:
  - label: ""
    status_equals: applied
`)
	_, err := labeler.LoadRules(path)
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestSaveRules_Roundtrip(t *testing.T) {
	rules := []labeler.Rule{
		{Label: "done", StatusEquals: "applied"},
		{Label: "stale", OlderThanDays: 14},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "out.yaml")
	if err := labeler.SaveRules(path, rules); err != nil {
		t.Fatalf("SaveRules: %v", err)
	}
	loaded, err := labeler.LoadRules(path)
	if err != nil {
		t.Fatalf("LoadRules: %v", err)
	}
	if len(loaded) != len(rules) {
		t.Fatalf("expected %d rules, got %d", len(rules), len(loaded))
	}
	if loaded[0].Label != "done" || loaded[1].OlderThanDays != 14 {
		t.Errorf("roundtrip mismatch: %+v", loaded)
	}
}

func TestLoadRules_EmptyFile(t *testing.T) {
	path := writeTempRules(t, "rules: []\n")
	rules, err := labeler.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}
