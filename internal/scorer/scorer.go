// Package scorer assigns a numeric priority score to patches based on
// configurable criteria such as tag weights, age, and retry count.
package scorer

import (
	"time"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/state"
)

// Result holds a patch name and its computed score.
type Result struct {
	PatchName string
	Score     float64
}

// Policy controls how scores are computed.
type Policy struct {
	// TagWeights maps tag names to additive score bonuses.
	TagWeights map[string]float64
	// AgeWeight is multiplied by the number of days since the patch was last
	// attempted (or zero if never attempted).
	AgeWeight float64
	// RetryPenalty is subtracted for each previous failed attempt.
	RetryPenalty float64
}

// DefaultPolicy returns a sensible default scoring policy.
func DefaultPolicy() Policy {
	return Policy{
		TagWeights:   map[string]float64{"critical": 10.0, "hotfix": 5.0},
		AgeWeight:    0.5,
		RetryPenalty: 2.0,
	}
}

// Score computes a Result for every patch in cfg using the supplied state and
// policy. Patches with no state entry are treated as never attempted.
func Score(cfg *config.Config, st *state.State, pol Policy) []Result {
	now := time.Now().UTC()
	results := make([]Result, 0, len(cfg.Patches))

	for _, p := range cfg.Patches {
		var score float64

		// Tag bonuses.
		for _, tag := range p.Tags {
			if w, ok := pol.TagWeights[tag]; ok {
				score += w
			}
		}

		// Age bonus and retry penalty from state.
		if entry, ok := st.Get(p.Name); ok {
			if !entry.AppliedAt.IsZero() {
				days := now.Sub(entry.AppliedAt).Hours() / 24
				score += pol.AgeWeight * days
			}
			score -= pol.RetryPenalty * float64(entry.Attempts)
		}

		results = append(results, Result{PatchName: p.Name, Score: score})
	}

	return results
}

// AnyNegative returns true if at least one result has a negative score.
func AnyNegative(results []Result) bool {
	for _, r := range results {
		if r.Score < 0 {
			return true
		}
	}
	return false
}
