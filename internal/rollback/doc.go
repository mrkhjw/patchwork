// Package rollback implements pre-patch backup and restoration for patchwork.
//
// Before applying a patch, callers should invoke Backup to snapshot the target
// file. If the patch needs to be reversed, Restore replaces the file with the
// snapshot and resets the patch state entry to StatusPending.
//
// Backups are stored under <repoPath>/.patchwork/backups/<patchName>.bak and
// are removed automatically after a successful Restore.
//
// Typical usage:
//
//	if err := rollback.Backup(repo, patch.Name, targetFile); err != nil { ... }
//	// apply patch ...
//	// later, to undo:
//	if err := rollback.Restore(repo, patch.Name, targetFile, st); err != nil { ... }
package rollback
