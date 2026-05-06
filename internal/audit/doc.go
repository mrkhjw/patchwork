// Package audit provides a structured, append-only event log for patchwork
// operations.
//
// Each time a patch is applied, skipped, rolled back, or fails, callers
// should append an Event to the Log.  The log can be persisted to disk as
// JSON and reloaded across runs, giving operators a full audit trail of
// every change patchwork has made to their repositories.
//
// Typical usage:
//
//	log, err := audit.Load(".patchwork/audit.json")
//	if err != nil { ... }
//
//	log.Append(patchName, repoPath, audit.EventApplied, "")
//
//	if err := log.Save(".patchwork/audit.json"); err != nil { ... }
//
// Use Filter to retrieve events of a specific kind for reporting or
// drift-detection purposes.
package audit
