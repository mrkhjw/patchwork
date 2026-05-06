// Package filter provides utilities for selecting subsets of patches
// and repos based on user-supplied criteria (tags, names, repo targets).
package filter

import (
	"github.com/yourorg/patchwork/internal/config"
)

// Options holds the criteria used to narrow down which patches to run.
type Options struct {
	// Names restricts processing to patches whose name is in this set.
	Names []string
	// Tags restricts processing to patches that carry at least one of these tags.
	Tags []string
	// Repos restricts processing to patches whose target repo is in this set.
	Repos []string
}

// Apply returns the subset of patches that match ALL non-empty criteria in opts.
func Apply(patches []config.Patch, opts Options) []config.Patch {
	nameSet := toSet(opts.Names)
	tagSet := toSet(opts.Tags)
	repoSet := toSet(opts.Repos)

	var out []config.Patch
	for _, p := range patches {
		if len(nameSet) > 0 && !nameSet[p.Name] {
			continue
		}
		if len(tagSet) > 0 && !hasAnyTag(p.Tags, tagSet) {
			continue
		}
		if len(repoSet) > 0 && !repoSet[p.Repo] {
			continue
		}
		out = append(out, p)
	}
	return out
}

// toSet converts a slice of strings into a lookup map.
func toSet(items []string) map[string]bool {
	if len(items) == 0 {
		return nil
	}
	m := make(map[string]bool, len(items))
	for _, v := range items {
		m[v] = true
	}
	return m
}

// hasAnyTag returns true if at least one of the patch's tags is in the set.
func hasAnyTag(patchTags []string, set map[string]bool) bool {
	for _, t := range patchTags {
		if set[t] {
			return true
		}
	}
	return false
}
