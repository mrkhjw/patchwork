// Package classifier categorises patches into risk tiers based on
// configurable rules so operators can prioritise review and rollout.
package classifier

import (
	"strings"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/state"
)

// Tier represents a risk classification level.
type Tier string

const (
	TierCritical Tier = "critical"
	TierHigh     Tier = "high"
	TierMedium   Tier = "medium"
	TierLow      Tier = "low"
)

// Result holds the classification outcome for a single patch.
type Result struct {
	PatchName string
	Repo      string
	Tier      Tier
	Reasons   []string
}

// Policy controls which conditions elevate risk.
type Policy struct {
	CriticalTags []string
	HighTags     []string
	FailedIsHigh bool
}

// DefaultPolicy returns a sensible out-of-the-box policy.
func DefaultPolicy() Policy {
	return Policy{
		CriticalTags: []string{"security", "breaking"},
		HighTags:     []string{"migration", "schema"},
		FailedIsHigh: true,
	}
}

// Classify assigns a risk tier to every patch in cfg, consulting st for
// current status information.
func Classify(cfg *config.Config, st *state.State, pol Policy) []Result {
	critSet := toSet(pol.CriticalTags)
	highSet := toSet(pol.HighTags)

	var results []Result
	for _, p := range cfg.Patches {
		var reasons []string
		tier := TierLow

		for _, tag := range p.Tags {
			low := strings.ToLower(tag)
			if critSet[low] {
				tier = TierCritical
				reasons = append(reasons, "critical tag: "+tag)
			} else if highSet[low] && tier != TierCritical {
				tier = TierHigh
				reasons = append(reasons, "high tag: "+tag)
			}
		}

		if entry, ok := st.Get(p.Name); ok {
			if entry.Status == "failed" && pol.FailedIsHigh && tier == TierLow {
				tier = TierHigh
				reasons = append(reasons, "status: failed")
			}
			if tier == TierLow && entry.Status == "pending" {
				tier = TierMedium
				reasons = append(reasons, "status: pending")
			}
		}

		results = append(results, Result{
			PatchName: p.Name,
			Repo:      p.Repo,
			Tier:      tier,
			Reasons:   reasons,
		})
	}
	return results
}

// AnyAbove returns true if any result has a tier strictly above threshold.
func AnyAbove(results []Result, threshold Tier) bool {
	order := map[Tier]int{TierLow: 0, TierMedium: 1, TierHigh: 2, TierCritical: 3}
	for _, r := range results {
		if order[r.Tier] > order[threshold] {
			return true
		}
	}
	return false
}

func toSet(tags []string) map[string]bool {
	s := make(map[string]bool, len(tags))
	for _, t := range tags {
		s[strings.ToLower(t)] = true
	}
	return s
}
