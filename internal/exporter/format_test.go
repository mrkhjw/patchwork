package exporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/exporter"
)

// TestExport_JSON_OmitsEmptyError verifies that entries without errors do not
// include the "error" key in JSON output (omitempty behaviour).
func TestExport_JSON_OmitsEmptyError(t *testing.T) {
	entries := []exporter.Entry{
		{PatchName: "clean", Repo: "svc", Status: "applied", AppliedAt: time.Now()},
	}
	var buf bytes.Buffer
	if err := exporter.Export(&buf, entries, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), `"error"`) {
		t.Error("expected 'error' key to be omitted when empty")
	}
}

// TestExport_CSV_ErrorColumn verifies the error column is present even when empty.
func TestExport_CSV_ErrorColumn(t *testing.T) {
	entries := []exporter.Entry{
		{PatchName: "ok", Repo: "svc", Status: "applied", AppliedAt: time.Now()},
	}
	var buf bytes.Buffer
	_ = exporter.Export(&buf, entries, exporter.FormatCSV)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 1 data row
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	// CSV row should have 5 fields
	fields := strings.Split(lines[1], ",")
	if len(fields) != 5 {
		t.Errorf("expected 5 CSV fields, got %d: %v", len(fields), fields)
	}
}

// TestExport_JSON_TimestampFormat verifies RFC3339 timestamp encoding.
func TestExport_JSON_TimestampFormat(t *testing.T) {
	stamp, _ := time.Parse(time.RFC3339, "2024-01-15T08:30:00Z")
	entries := []exporter.Entry{
		{PatchName: "ts-test", Repo: "svc", Status: "applied", AppliedAt: stamp},
	}
	var buf bytes.Buffer
	_ = exporter.Export(&buf, entries, exporter.FormatJSON)
	if !strings.Contains(buf.String(), "2024-01-15T08:30:00Z") {
		t.Errorf("expected RFC3339 timestamp in JSON output, got:\n%s", buf.String())
	}
}
