// Package cleaner removes patch state entries that have been applied
// successfully and whose patch files no longer exist on disk.
package cleaner

import (
	"os"

	"patchwork/internal/config"
	"patchwork/internal/state"
)

// Result describes the outcome for a single cleaned entry.
type Result struct {
	PatchName string
	Repo      string
	Reason    string
	Removed   bool
}

// Clean inspects the current state against the provided config and removes
// entries whose patch file is gone and whose status is "applied". Entries
// that are pending or failed are left untouched so they can be retried.
func Clean(cfg *config.Config, st *state.State) []Result {
	var results []Result

	known := make(map[string]bool, len(cfg.Patches))
	for _, p := range cfg.Patches {
		known[p.Name] = true
	}

	for _, entry := range st.All() {
		result := Result{
			PatchName: entry.PatchName,
			Repo:      entry.Repo,
		}

		if !known[entry.PatchName] {
			if entry.Status == "applied" {
				st.Remove(entry.PatchName, entry.Repo)
				result.Reason = "patch removed from config"
				result.Removed = true
			} else {
				result.Reason = "patch not in config but status is not applied — skipped"
			}
			results = append(results, result)
			continue
		}

		patch := patchByName(cfg, entry.PatchName)
		if patch == nil {
			continue
		}

		if _, err := os.Stat(patch.File); os.IsNotExist(err) && entry.Status == "applied" {
			st.Remove(entry.PatchName, entry.Repo)
			result.Reason = "patch file missing from disk"
			result.Removed = true
			results = append(results, result)
		}
	}

	return results
}

// AnyRemoved returns true if at least one entry was removed.
func AnyRemoved(results []Result) bool {
	for _, r := range results {
		if r.Removed {
			return true
		}
	}
	return false
}

func patchByName(cfg *config.Config, name string) *config.Patch {
	for i := range cfg.Patches {
		if cfg.Patches[i].Name == name {
			return &cfg.Patches[i]
		}
	}
	return nil
}
