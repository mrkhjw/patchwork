package templater

import (
	"fmt"
	"strings"
)

// Format returns a human-readable summary of all render results.
func Format(results []Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s  %-6s  %s\n", "PATCH", "OK", "OUTPUT / ERROR"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, r := range results {
		ok := "yes"
		detail := truncate(r.Output, 40)
		if r.Err != nil {
			ok = "no"
			detail = r.Err.Error()
		}
		sb.WriteString(fmt.Sprintf("%-20s  %-6s  %s\n", r.PatchName, ok, detail))
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
