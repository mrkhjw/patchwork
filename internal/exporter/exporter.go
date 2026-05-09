// Package exporter serialises patch run results to common interchange formats
// (JSON and CSV) so external tools can consume patchwork output.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Entry represents a single patch result record for export.
type Entry struct {
	PatchName string    `json:"patch_name"`
	Repo      string    `json:"repo"`
	Status    string    `json:"status"`
	AppliedAt time.Time `json:"applied_at"`
	Error     string    `json:"error,omitempty"`
}

// Format enumerates supported export formats.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Export writes entries to w in the requested format.
// Returns an error if the format is unsupported or writing fails.
func Export(w io.Writer, entries []Entry, format Format) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, entries)
	case FormatCSV:
		return exportCSV(w, entries)
	default:
		return fmt.Errorf("exporter: unsupported format %q", format)
	}
}

func exportJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func exportCSV(w io.Writer, entries []Entry) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"patch_name", "repo", "status", "applied_at", "error"}); err != nil {
		return fmt.Errorf("exporter: write csv header: %w", err)
	}
	for _, e := range entries {
		row := []string{
			e.PatchName,
			e.Repo,
			e.Status,
			e.AppliedAt.Format(time.RFC3339),
			e.Error,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("exporter: write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

// AnyFailed returns true if any entry carries a non-empty error string.
func AnyFailed(entries []Entry) bool {
	for _, e := range entries {
		if e.Error != "" {
			return true
		}
	}
	return false
}
