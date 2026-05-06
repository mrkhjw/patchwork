// Package rollback provides functionality to reverse previously applied patches
// by recording pre-patch state and restoring it on demand.
package rollback

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patchwork/internal/state"
)

// BackupDir is the subdirectory within the repo used to store patch backups.
const BackupDir = ".patchwork/backups"

// Backup captures the current content of targetFile before a patch is applied.
// It stores the backup keyed by patchName under repoPath.
func Backup(repoPath, patchName, targetFile string) error {
	backupPath := backupFilePath(repoPath, patchName)
	if err := os.MkdirAll(filepath.Dir(backupPath), 0o755); err != nil {
		return fmt.Errorf("rollback: mkdir: %w", err)
	}

	data, err := os.ReadFile(targetFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("rollback: read original: %w", err)
	}
	// If file did not exist, store empty sentinel.
	if os.IsNotExist(err) {
		data = []byte{}
	}

	if err := os.WriteFile(backupPath, data, 0o644); err != nil {
		return fmt.Errorf("rollback: write backup: %w", err)
	}
	return nil
}

// Restore reverses a previously applied patch by restoring the backed-up file
// and updating the patch state entry to StatusPending.
func Restore(repoPath, patchName, targetFile string, st *state.State) error {
	backupPath := backupFilePath(repoPath, patchName)

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("rollback: read backup: %w", err)
	}

	if len(data) == 0 {
		// Original file did not exist; remove the patched file.
		if err := os.Remove(targetFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rollback: remove patched file: %w", err)
		}
	} else {
		if err := os.WriteFile(targetFile, data, 0o644); err != nil {
			return fmt.Errorf("rollback: restore file: %w", err)
		}
	}

	_ = os.Remove(backupPath)

	entry, ok := st.Get(patchName)
	if ok {
		entry.Status = state.StatusPending
		entry.AppliedAt = ""
		st.Upsert(entry)
	}
	return nil
}

// HasBackup reports whether a backup exists for the given patch in repoPath.
func HasBackup(repoPath, patchName string) bool {
	_, err := os.Stat(backupFilePath(repoPath, patchName))
	return err == nil
}

func backupFilePath(repoPath, patchName string) string {
	return filepath.Join(repoPath, BackupDir, patchName+".bak")
}
