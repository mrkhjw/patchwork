// Package summarizer aggregates patch run results into a structured summary
// suitable for reporting, notifications, and audit logging.
package summarizer

import (
	"fmt"
	"time"

	"patchwork/internal/state"
)

// PatchSummary holds aggregated statistics for a single patch run.
type PatchSummary struct {
	PatchName  string
	Repo       string
	Status     string
	AppliedAt  time.Time
	DurationMs int64
}

// RunSummary holds the full aggregated result of a patchwork run.
type RunSummary struct {
	Total    int
	Applied  int
	Skipped  int
	Failed   int
	Patches  []PatchSummary
	RunAt    time.Time
}

// Build constructs a RunSummary from a slice of state entries.
func Build(entries []state.Entry) RunSummary {
	summary := RunSummary{
		RunAt:   time.Now().UTC(),
		Total:   len(entries),
		Patches: make([]PatchSummary, 0, len(entries)),
	}

	for _, e := range entries {
		ps := PatchSummary{
			PatchName: e.PatchName,
			Repo:      e.Repo,
			Status:    e.Status,
			AppliedAt: e.AppliedAt,
		}
		summary.Patches = append(summary.Patches, ps)

		switch e.Status {
		case "applied":
			summary.Applied++
		case "skipped":
			summary.Skipped++
		case "failed":
			summary.Failed++
		}
	}

	return summary
}

// Format returns a human-readable multi-line summary string.
func Format(s RunSummary) string {
	return fmt.Sprintf(
		"Run at: %s\nTotal: %d | Applied: %d | Skipped: %d | Failed: %d",
		s.RunAt.Format(time.RFC3339),
		s.Total,
		s.Applied,
		s.Skipped,
		s.Failed,
	)
}

// AnyFailed returns true if at least one patch entry has a failed status.
func AnyFailed(s RunSummary) bool {
	return s.Failed > 0
}
