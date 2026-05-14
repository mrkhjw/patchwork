package trimmer

import (
	"fmt"
	"strings"
)

// Format returns a human-readable summary of trimmed entries.
func Format(results []Result) string {
	if len(results) == 0 {
		return "trimmer: nothing removed\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "trimmer: removed %d entr", len(results))
	if len(results) == 1 {
		sb.WriteString("y\n")
	} else {
		sb.WriteString("ies\n")
	}

	for _, r := range results {
		fmt.Fprintf(&sb, "  - %-30s  reason: %s\n", r.PatchName, r.Reason)
	}
	return sb.String()
}
