// Package tagger provides utilities for querying and summarising patch tags
// across a set of patch definitions.
package tagger

import (
	"sort"

	"github.com/yourorg/patchwork/internal/config"
)

// TagIndex maps each tag to the patch names that carry it.
type TagIndex map[string][]string

// Build constructs a TagIndex from a slice of patch configurations.
// Patches with no tags are silently skipped.
func Build(patches []config.Patch) TagIndex {
	idx := make(TagIndex)
	for _, p := range patches {
		for _, tag := range p.Tags {
			idx[tag] = append(idx[tag], p.Name)
		}
	}
	// Sort each slice so output is deterministic.
	for tag := range idx {
		sort.Strings(idx[tag])
	}
	return idx
}

// AllTags returns a sorted list of every unique tag present in the index.
func (idx TagIndex) AllTags() []string {
	tags := make([]string, 0, len(idx))
	for t := range idx {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}

// PatchesForTag returns the patch names associated with the given tag.
// Returns nil if the tag is not present in the index.
func (idx TagIndex) PatchesForTag(tag string) []string {
	return idx[tag]
}

// TagsForPatch returns a sorted list of tags that are applied to the named
// patch. Returns nil when the patch has no tags or is not found.
func TagsForPatch(patches []config.Patch, name string) []string {
	for _, p := range patches {
		if p.Name == name {
			if len(p.Tags) == 0 {
				return nil
			}
			out := make([]string, len(p.Tags))
			copy(out, p.Tags)
			sort.Strings(out)
			return out
		}
	}
	return nil
}
