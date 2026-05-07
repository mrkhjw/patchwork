// Package pruner removes stale state entries for patches that no longer exist
// in the configuration, keeping the state file clean over time.
package pruner

import (
	"fmt"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/state"
)

// Result holds the outcome of a prune operation for a single entry.
type Result struct {
	PatchName string
	Removed   bool
	Reason    string
}

// Prune removes state entries whose patch names are not present in cfg.
// It returns one Result per stale entry that was removed.
func Prune(st *state.State, cfg *config.Config) ([]Result, error) {
	known := make(map[string]struct{}, len(cfg.Patches))
	for _, p := range cfg.Patches {
		known[p.Name] = struct{}{}
	}

	all := st.All()
	var results []Result

	for _, entry := range all {
		if _, ok := known[entry.PatchName]; ok {
			continue
		}
		if err := st.Remove(entry.PatchName); err != nil {
			return results, fmt.Errorf("pruner: remove %q: %w", entry.PatchName, err)
		}
		results = append(results, Result{
			PatchName: entry.PatchName,
			Removed:   true,
			Reason:    "patch no longer in config",
		})
	}

	return results, nil
}

// AnyRemoved returns true if at least one entry was pruned.
func AnyRemoved(results []Result) bool {
	for _, r := range results {
		if r.Removed {
			return true
		}
	}
	return false
}
