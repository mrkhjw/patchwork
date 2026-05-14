// Package trimmer removes patch state entries that exceed a maximum age or count threshold.
package trimmer

import (
	"fmt"
	"sort"
	"time"

	"github.com/patchwork/internal/state"
)

// Policy controls how trimming is performed.
type Policy struct {
	// MaxAge removes entries whose last-updated timestamp is older than this duration.
	// Zero means no age-based trimming.
	MaxAge time.Duration

	// MaxEntries keeps only the N most-recently updated entries.
	// Zero means no count-based trimming.
	MaxEntries int
}

// Result describes a single trimmed entry.
type Result struct {
	PatchName string
	Reason    string
}

// Trim removes stale entries from s according to p and returns the list of
// removed entries. The caller is responsible for persisting the mutated state.
func Trim(s *state.State, p Policy) []Result {
	var removed []Result

	entries := s.All()
	now := time.Now()

	// Age-based pass.
	if p.MaxAge > 0 {
		for _, e := range entries {
			if !e.UpdatedAt.IsZero() && now.Sub(e.UpdatedAt) > p.MaxAge {
				s.Delete(e.PatchName)
				removed = append(removed, Result{
					PatchName: e.PatchName,
					Reason:    fmt.Sprintf("older than %s", p.MaxAge),
				})
			}
		}
		// Refresh slice after deletions.
		entries = s.All()
	}

	// Count-based pass: keep the N most-recently updated.
	if p.MaxEntries > 0 && len(entries) > p.MaxEntries {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].UpdatedAt.After(entries[j].UpdatedAt)
		})
		for _, e := range entries[p.MaxEntries:] {
			s.Delete(e.PatchName)
			removed = append(removed, Result{
				PatchName: e.PatchName,
				Reason:    "exceeds max entry count",
			})
		}
	}

	return removed
}

// AnyRemoved returns true when results contains at least one entry.
func AnyRemoved(results []Result) bool {
	return len(results) > 0
}
