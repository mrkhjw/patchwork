// Package trimmer provides age- and count-based eviction of patch state entries.
//
// Use Trim to remove entries from a [state.State] that are either too old or
// exceed a maximum count. The caller controls persistence — Trim only mutates
// the in-memory state.
//
// Example:
//
//	p := trimmer.Policy{MaxAge: 30 * 24 * time.Hour, MaxEntries: 500}
//	removed := trimmer.Trim(s, p)
//	if trimmer.AnyRemoved(removed) {
//		fmt.Print(trimmer.Format(removed))
//		_ = s.Save(statePath)
//	}
package trimmer
