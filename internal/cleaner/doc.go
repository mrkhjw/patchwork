// Package cleaner provides utilities for removing stale patch state entries
// from patchwork's state store.
//
// An entry is considered stale when:
//   - Its status is "applied" AND the patch file no longer exists on disk, or
//   - Its patch name is no longer present in the active configuration.
//
// Entries with status "pending" or "failed" are never removed automatically,
// as they may still need attention or retrying.
//
// # Staleness Rules
//
// The cleaner distinguishes between two kinds of staleness:
//
//   - File-missing: the patch was applied but its source file has since been
//     deleted or moved. The state entry is safe to drop because the patch
//     cannot be re-applied or rolled back without the file.
//
//   - Config-missing: the patch name no longer appears in the active
//     configuration. This typically happens when a patch is intentionally
//     retired and removed from the config file.
//
// # Typical Usage
//
//	results := cleaner.Clean(cfg, st)
//	if cleaner.AnyRemoved(results) {
//		_ = st.Save(statePath)
//	}
package cleaner
