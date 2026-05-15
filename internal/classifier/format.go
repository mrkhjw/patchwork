package classifier

import (
	"fmt"
	"strings"
)

// Format returns a human-readable summary table of classification results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no patches to classify\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-25s %-10s %s\n",
		"PATCH", "REPO", "TIER", "REASONS"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")

	for _, r := range results {
		reasons := "-"
		if len(r.Reasons) > 0 {
			reasons = strings.Join(r.Reasons, "; ")
		}
		sb.WriteString(fmt.Sprintf("%-30s %-25s %-10s %s\n",
			truncate(r.PatchName, 29),
			truncate(r.Repo, 24),
			string(r.Tier),
			reasons,
		))
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
