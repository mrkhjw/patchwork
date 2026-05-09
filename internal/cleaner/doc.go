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
// Typical usage:
//
//	results := cleaner.Clean(cfg, st)
//	if cleaner.AnyRemoved(results) {
//		_ = st.Save(statePath)
//	}
package cleaner
