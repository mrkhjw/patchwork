// Package reporter provides formatted output for patch status summaries.
package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/yourorg/patchwork/internal/state"
)

// Status constants mirror state values for display purposes.
const (
	StatusApplied = "applied"
	StatusPending = "pending"
	StatusFailed  = "failed"
)

// Row represents a single line in the status report.
type Row struct {
	Repo      string
	Patch     string
	Status    string
	AppliedAt string
}

// Report writes a tabular summary of patch states to w.
func Report(w io.Writer, entries []state.Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "REPO\tPATCH\tSTATUS\tAPPLIED AT")
	fmt.Fprintln(tw, "----\t-----\t------\t----------")
	for _, e := range entries {
		appliedAt := "-"
		if !e.AppliedAt.IsZero() {
			appliedAt = e.AppliedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.Repo, e.Patch, e.Status, appliedAt)
	}
	return tw.Flush()
}

// ReportToStdout is a convenience wrapper that writes to os.Stdout.
func ReportToStdout(entries []state.Entry) error {
	return Report(os.Stdout, entries)
}

// Summary prints a brief count summary to w.
func Summary(w io.Writer, entries []state.Entry) {
	counts := map[string]int{}
	for _, e := range entries {
		counts[e.Status]++
	}
	fmt.Fprintf(w, "Summary: %d applied, %d pending, %d failed\n",
		counts[StatusApplied],
		counts[StatusPending],
		counts[StatusFailed],
	)
}
