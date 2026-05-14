// Package inspector provides a summary view of patch health across repos,
// combining state, config, and digest information into a single report.
package inspector

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/state"
)

// Result holds the inspection output for a single patch.
type Result struct {
	Name      string
	Repo      string
	Status    string
	AppliedAt time.Time
	Tags      []string
	Warnings  []string
}

// AnyWarnings returns true if at least one result carries a warning.
func AnyWarnings(results []Result) bool {
	for _, r := range results {
		if len(r.Warnings) > 0 {
			return true
		}
	}
	return false
}

// Inspect cross-references cfg patches with st entries and returns one Result
// per patch, annotating any anomalies as warnings.
func Inspect(cfg *config.Config, st *state.State) []Result {
	results := make([]Result, 0, len(cfg.Patches))

	for _, p := range cfg.Patches {
		entry, found := st.Get(p.Name)

		r := Result{
			Name: p.Name,
			Repo: p.Repo,
			Tags: p.Tags,
		}

		if !found {
			r.Status = "pending"
			r.Warnings = append(r.Warnings, "no state entry found")
		} else {
			r.Status = entry.Status
			r.AppliedAt = entry.AppliedAt

			if entry.Status == "applied" && entry.AppliedAt.IsZero() {
				r.Warnings = append(r.Warnings, "applied status but zero timestamp")
			}
			if entry.Status == "failed" {
				r.Warnings = append(r.Warnings, "patch is in failed state")
			}
		}

		if len(p.Tags) == 0 {
			r.Warnings = append(r.Warnings, "patch has no tags")
		}

		results = append(results, r)
	}

	return results
}

// Format renders inspection results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no patches to inspect\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-30s %-20s %-10s %s\n", "NAME", "REPO", "STATUS", "WARNINGS")
	fmt.Fprintln(&sb, strings.Repeat("-", 80))

	for _, r := range results {
		warn := strings.Join(r.Warnings, "; ")
		if warn == "" {
			warn = "-"
		}
		fmt.Fprintf(&sb, "%-30s %-20s %-10s %s\n", r.Name, r.Repo, r.Status, warn)
	}

	return sb.String()
}
