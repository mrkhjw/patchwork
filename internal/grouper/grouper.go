// Package grouper groups patches by a chosen dimension (repo, tag, or status)
// and returns an ordered map of bucket name → patch slice.
package grouper

import (
	"fmt"
	"sort"

	"github.com/yourorg/patchwork/internal/config"
)

// By is the grouping dimension.
type By string

const (
	ByRepo   By = "repo"
	ByTag    By = "tag"
	ByStatus By = "status"
)

// Result holds the ordered bucket keys and the grouped patches.
type Result struct {
	Keys    []string
	Buckets map[string][]config.Patch
}

// Group partitions patches according to the requested dimension.
// statuses maps patch name → status string (e.g. "applied", "pending").
func Group(patches []config.Patch, by By, statuses map[string]string) (*Result, error) {
	buckets := make(map[string][]config.Patch)

	for _, p := range patches {
		var keys []string

		switch by {
		case ByRepo:
			keys = p.Repos
			if len(keys) == 0 {
				keys = []string{"(no repo)"}
			}
		case ByTag:
			keys = p.Tags
			if len(keys) == 0 {
				keys = []string{"(untagged)"}
			}
		case ByStatus:
			status, ok := statuses[p.Name]
			if !ok {
				status = "unknown"
			}
			keys = []string{status}
		default:
			return nil, fmt.Errorf("grouper: unknown dimension %q", by)
		}

		for _, k := range keys {
			buckets[k] = append(buckets[k], p)
		}
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &Result{Keys: keys, Buckets: buckets}, nil
}

// Counts returns a map of bucket name → patch count for quick summaries.
func (r *Result) Counts() map[string]int {
	out := make(map[string]int, len(r.Keys))
	for _, k := range r.Keys {
		out[k] = len(r.Buckets[k])
	}
	return out
}
