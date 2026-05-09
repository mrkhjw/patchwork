// Package sorter provides utilities for ordering patches by various criteria
// such as name, repo, status, or score.
package sorter

import (
	"sort"
	"strings"

	"patchwork/internal/config"
)

// Dimension controls the field used for sorting.
type Dimension string

const (
	ByName   Dimension = "name"
	ByRepo   Dimension = "repo"
	ByStatus Dimension = "status"
)

// Entry holds a patch alongside its current status for sorting purposes.
type Entry struct {
	Patch  config.Patch
	Status string
}

// Result is the ordered list of entries after sorting.
type Result struct {
	Entries []Entry
	Dim     Dimension
}

// Sort orders the given entries by the requested dimension.
// Unknown dimensions fall back to sorting by name.
// Sorting is always stable and case-insensitive.
func Sort(entries []Entry, dim Dimension) Result {
	copied := make([]Entry, len(entries))
	copy(copied, entries)

	sort.SliceStable(copied, func(i, j int) bool {
		switch dim {
		case ByRepo:
			ri := strings.ToLower(copied[i].Patch.Repo)
			rj := strings.ToLower(copied[j].Patch.Repo)
			if ri != rj {
				return ri < rj
			}
			// secondary sort by name for stability
			return strings.ToLower(copied[i].Patch.Name) < strings.ToLower(copied[j].Patch.Name)
		case ByStatus:
			si := strings.ToLower(copied[i].Status)
			sj := strings.ToLower(copied[j].Status)
			if si != sj {
				return si < sj
			}
			return strings.ToLower(copied[i].Patch.Name) < strings.ToLower(copied[j].Patch.Name)
		default: // ByName
			return strings.ToLower(copied[i].Patch.Name) < strings.ToLower(copied[j].Patch.Name)
		}
	})

	return Result{Entries: copied, Dim: dim}
}

// Names returns the ordered patch names from a Result.
func (r Result) Names() []string {
	out := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		out[i] = e.Patch.Name
	}
	return out
}
