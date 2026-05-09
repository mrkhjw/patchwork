// Package matcher provides fuzzy and exact matching utilities for patch names,
// repos, and tags to support interactive search and filtering workflows.
package matcher

import (
	"strings"

	"github.com/patchwork/internal/config"
)

// Result holds a matched patch along with its relevance score.
type Result struct {
	Patch config.Patch
	Score int
}

// Match returns patches whose name, repo, or tags contain the query string
// (case-insensitive). Results are ordered: exact name match first, then
// partial name, then repo/tag matches.
func Match(patches []config.Patch, query string) []Result {
	if query == "" {
		return allAsResults(patches)
	}

	q := strings.ToLower(query)
	var results []Result

	for _, p := range patches {
		score := scoreMatch(p, q)
		if score > 0 {
			results = append(results, Result{Patch: p, Score: score})
		}
	}

	sortByScore(results)
	return results
}

// AnyMatch returns true if at least one patch matches the query.
func AnyMatch(patches []config.Patch, query string) bool {
	return len(Match(patches, query)) > 0
}

func scoreMatch(p config.Patch, q string) int {
	name := strings.ToLower(p.Name)
	if name == q {
		return 100
	}
	if strings.HasPrefix(name, q) {
		return 80
	}
	if strings.Contains(name, q) {
		return 60
	}
	if strings.Contains(strings.ToLower(p.Repo), q) {
		return 40
	}
	for _, tag := range p.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return 20
		}
	}
	return 0
}

func allAsResults(patches []config.Patch) []Result {
	out := make([]Result, len(patches))
	for i, p := range patches {
		out[i] = Result{Patch: p, Score: 0}
	}
	return out
}

func sortByScore(results []Result) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Score > results[j-1].Score; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
