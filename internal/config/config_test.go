package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "patchwork-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `version: "1"
repos:
  - name: core
    path: /projects/core
    remote: origin
  - name: api
    path: /projects/api
`
	path := writeTemp(t, raw)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(cfg.Repos))
	}
	if cfg.Repos[0].Name != "core" {
		t.Errorf("expected first repo name 'core', got %q", cfg.Repos[0].Name)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_DuplicateName(t *testing.T) {
	raw := `version: "1"
repos:
  - name: core
    path: /projects/core
  - name: core
    path: /projects/core2
`
	path := writeTemp(t, raw)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate repo name")
	}
}

func TestSaveRoundtrip(t *testing.T) {
	cfg := &config.Config{
		Version: "1",
		Repos: []config.Repo{
			{Name: "myrepo", Path: "/tmp/myrepo", Remote: "origin"},
		},
	}
	dest := filepath.Join(t.TempDir(), "out.yaml")
	if err := cfg.Save(dest); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := config.Load(dest)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if loaded.Repos[0].Name != "myrepo" {
		t.Errorf("roundtrip mismatch: got %q", loaded.Repos[0].Name)
	}
}
