package watchdog

import (
	"fmt"
	"strings"
)

// Format renders a human-readable report of watchdog alerts.
func Format(alerts []Alert) string {
	if len(alerts) == 0 {
		return "watchdog: no alerts\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "watchdog: %d alert(s)\n", len(alerts))
	fmt.Fprintf(&sb, "%-12s %-20s %-12s %-10s %s\n", "SEVERITY", "PATCH", "REPO", "STATUS", "AGE")
	fmt.Fprintln(&sb, strings.Repeat("-", 72))
	for _, a := range alerts {
		fmt.Fprintf(&sb, "%-12s %-20s %-12s %-10s %s\n",
			string(a.Severity),
			truncate(a.PatchName, 20),
			truncate(a.Repo, 12),
			a.Status,
			a.Age,
		)
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
