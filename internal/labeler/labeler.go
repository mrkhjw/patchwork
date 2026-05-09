// Package labeler assigns computed labels to patches based on configurable rules.
// Labels differ from tags in that they are derived automatically from patch
// metadata (e.g. repo, status, age) rather than declared by the author.
package labeler

import (
	"fmt"
	"strings"
	"time"
)

// Rule describes a single labeling rule.
type Rule struct {
	// Label is the string attached to matching patches.
	Label string
	// RepoContains matches when the repo path contains the given substring.
	RepoContains string
	// OlderThanDays matches when the patch was applied more than N days ago.
	OlderThanDays int
	// StatusEquals matches when the patch status equals the given string.
	StatusEquals string
}

// PatchMeta holds the fields labeler inspects.
type PatchMeta struct {
	Name      string
	Repo      string
	Status    string
	AppliedAt time.Time
}

// Result holds the labels assigned to a single patch.
type Result struct {
	Patch  PatchMeta
	Labels []string
}

// AnyFailed returns true when at least one result carries the "error" label.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		for _, l := range r.Labels {
			if l == "error" {
				return true
			}
		}
	}
	return false
}

// Label applies all rules to every patch and returns the results.
func Label(patches []PatchMeta, rules []Rule) []Result {
	results := make([]Result, 0, len(patches))
	for _, p := range patches {
		results = append(results, Result{
			Patch:  p,
			Labels: applyRules(p, rules),
		})
	}
	return results
}

// Format returns a human-readable summary line for a result.
func Format(r Result) string {
	if len(r.Labels) == 0 {
		return fmt.Sprintf("%s: (no labels)", r.Patch.Name)
	}
	return fmt.Sprintf("%s: [%s]", r.Patch.Name, strings.Join(r.Labels, ", "))
}

func applyRules(p PatchMeta, rules []Rule) []string {
	seen := map[string]struct{}{}
	var labels []string
	for _, r := range rules {
		if r.Label == "" {
			continue
		}
		if matches(p, r) {
			if _, ok := seen[r.Label]; !ok {
				seen[r.Label] = struct{}{}
				labels = append(labels, r.Label)
			}
		}
	}
	return labels
}

func matches(p PatchMeta, r Rule) bool {
	if r.RepoContains != "" && !strings.Contains(p.Repo, r.RepoContains) {
		return false
	}
	if r.StatusEquals != "" && p.Status != r.StatusEquals {
		return false
	}
	if r.OlderThanDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -r.OlderThanDays)
		if !p.AppliedAt.Before(cutoff) {
			return false
		}
	}
	return true
}
