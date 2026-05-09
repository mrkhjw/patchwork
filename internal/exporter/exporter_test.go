package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/exporter"
)

func sampleEntries() []exporter.Entry {
	t0, _ := time.Parse(time.RFC3339, "2024-06-01T12:00:00Z")
	return []exporter.Entry{
		{PatchName: "add-header", Repo: "api", Status: "applied", AppliedAt: t0},
		{PatchName: "fix-env", Repo: "worker", Status: "failed", AppliedAt: t0, Error: "patch rejected"},
	}
}

func TestExport_JSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, sampleEntries(), exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []exporter.Entry
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].PatchName != "add-header" {
		t.Errorf("unexpected patch name: %s", out[0].PatchName)
	}
}

func TestExport_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, sampleEntries(), exporter.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines (header+2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "patch_name") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExport_CSV_ContainsData(t *testing.T) {
	var buf bytes.Buffer
	_ = exporter.Export(&buf, sampleEntries(), exporter.FormatCSV)
	if !strings.Contains(buf.String(), "fix-env") {
		t.Error("expected 'fix-env' in CSV output")
	}
	if !strings.Contains(buf.String(), "patch rejected") {
		t.Error("expected error message in CSV output")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, sampleEntries(), exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnyFailed_True(t *testing.T) {
	if !exporter.AnyFailed(sampleEntries()) {
		t.Error("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	entries := []exporter.Entry{
		{PatchName: "ok", Repo: "api", Status: "applied"},
	}
	if exporter.AnyFailed(entries) {
		t.Error("expected AnyFailed to return false")
	}
}

func TestExport_EmptyEntries_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, nil, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty JSON output for nil entries")
	}
}
