package scorer

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// SortDesc sorts results in descending score order (highest priority first).
func SortDesc(results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
}

// Format writes a human-readable score table to w.
// Results are printed in the order supplied; call SortDesc first if desired.
func Format(w io.Writer, results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "(no patches to score)")
		return err
	}

	header := fmt.Sprintf("%-30s %10s", "PATCH", "SCORE")
	sep := strings.Repeat("-", len(header))

	if _, err := fmt.Fprintln(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, sep); err != nil {
		return err
	}

	for _, r := range results {
		line := fmt.Sprintf("%-30s %10.2f", r.PatchName, r.Score)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
