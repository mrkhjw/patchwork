// Package differ compares two snapshots and produces a human-readable drift report.
package differ

import (
	"fmt"
	"strings"

	"github.com/patchwork/internal/snapshot"
)

// Entry describes a single drift item between two snapshots.
type Entry struct {
	Repo   string
	Patch  string
	Old    string
	New    string
	Drifted bool
}

// Compare takes a baseline and current snapshot and returns drift entries.
func Compare(baseline, current *snapshot.Snapshot) []Entry {
	var entries []Entry

	for key, cur := range current.Entries {
		base, ok := baseline.Entries[key]
		if !ok {
			entries = append(entries, Entry{
				Repo:    cur.Repo,
				Patch:   cur.Patch,
				Old:     "(absent)",
				New:     cur.Status,
				Drifted: true,
			})
			continue
		}
		drifted := base.Status != cur.Status
		entries = append(entries, Entry{
			Repo:    cur.Repo,
			Patch:   cur.Patch,
			Old:     base.Status,
			New:     cur.Status,
			Drifted: drifted,
		})
	}

	// Detect entries present in baseline but missing in current.
	for key, base := range baseline.Entries {
		if _, ok := current.Entries[key]; !ok {
			entries = append(entries, Entry{
				Repo:    base.Repo,
				Patch:   base.Patch,
				Old:     base.Status,
				New:     "(absent)",
				Drifted: true,
			})
		}
	}

	return entries
}

// Format renders drift entries as a plain-text report string.
func Format(entries []Entry) string {
	if len(entries) == 0 {
		return "No drift detected.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-30s %-12s %-12s %s\n",
		"REPO", "PATCH", "OLD", "NEW", "DRIFTED"))
	sb.WriteString(strings.Repeat("-", 95) + "\n")
	for _, e := range entries {
		drifted := "no"
		if e.Drifted {
			drifted = "YES"
		}
		sb.WriteString(fmt.Sprintf("%-30s %-30s %-12s %-12s %s\n",
			e.Repo, e.Patch, e.Old, e.New, drifted))
	}
	return sb.String()
}

// AnyDrifted returns true if at least one entry has drifted.
func AnyDrifted(entries []Entry) bool {
	for _, e := range entries {
		if e.Drifted {
			return true
		}
	}
	return false
}
