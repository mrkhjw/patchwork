package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/reporter"
	"github.com/yourorg/patchwork/internal/state"
)

func sampleEntries() []state.Entry {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return []state.Entry{
		{Repo: "repo-a", Patch: "fix-001.patch", Status: "applied", AppliedAt: now},
		{Repo: "repo-b", Patch: "fix-002.patch", Status: "pending"},
		{Repo: "repo-c", Patch: "fix-003.patch", Status: "failed"},
	}
}

func TestReport_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Report(&buf, sampleEntries()); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	out := buf.String()
	for _, header := range []string{"REPO", "PATCH", "STATUS", "APPLIED AT"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestReport_ContainsEntries(t *testing.T) {
	var buf bytes.Buffer
	_ = reporter.Report(&buf, sampleEntries())
	out := buf.String()
	for _, want := range []string{"repo-a", "fix-002.patch", "failed", "2024-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output\ngot:\n%s", want, out)
		}
	}
}

func TestReport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Report(&buf, []state.Entry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "REPO") {
		t.Error("expected header even for empty entries")
	}
}

func TestSummary_Counts(t *testing.T) {
	var buf bytes.Buffer
	reporter.Summary(&buf, sampleEntries())
	out := buf.String()
	if !strings.Contains(out, "1 applied") {
		t.Errorf("expected '1 applied' in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 pending") {
		t.Errorf("expected '1 pending' in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 failed") {
		t.Errorf("expected '1 failed' in summary, got: %s", out)
	}
}

func TestReport_ZeroAppliedAt(t *testing.T) {
	var buf bytes.Buffer
	entries := []state.Entry{
		{Repo: "repo-x", Patch: "p.patch", Status: "pending"},
	}
	_ = reporter.Report(&buf, entries)
	if !strings.Contains(buf.String(), "-") {
		t.Error("expected '-' for zero AppliedAt")
	}
}
