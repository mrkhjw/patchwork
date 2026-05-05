package patch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Status represents the application state of a patch.
type Status string

const (
	StatusPending  Status = "pending"
	StatusApplied  Status = "applied"
	StatusConflict Status = "conflict"
	StatusUnknown  Status = "unknown"
)

// Patch describes a single patch file and its target repository.
type Patch struct {
	Name    string
	File    string
	RepoDir string
}

// Apply runs `git apply` for the patch in the target repository.
// It returns an error if the patch cannot be applied cleanly.
func Apply(p Patch) error {
	if _, err := os.Stat(p.File); err != nil {
		return fmt.Errorf("patch file not found: %s", p.File)
	}

	absFile, err := filepath.Abs(p.File)
	if err != nil {
		return fmt.Errorf("resolving patch path: %w", err)
	}

	cmd := exec.Command("git", "apply", "--check", absFile)
	cmd.Dir = p.RepoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("patch would not apply cleanly: %w\n%s", err, string(out))
	}

	cmd = exec.Command("git", "apply", absFile)
	cmd.Dir = p.RepoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("applying patch: %w\n%s", err, string(out))
	}

	return nil
}

// Check returns the Status of a patch against the target repository
// without modifying the working tree.
func Check(p Patch) Status {
	if _, err := os.Stat(p.File); err != nil {
		return StatusUnknown
	}

	absFile, err := filepath.Abs(p.File)
	if err != nil {
		return StatusUnknown
	}

	// --reverse --check detects whether the patch is already applied.
	cmd := exec.Command("git", "apply", "--reverse", "--check", absFile)
	cmd.Dir = p.RepoDir
	if err := cmd.Run(); err == nil {
		return StatusApplied
	}

	cmd = exec.Command("git", "apply", "--check", absFile)
	cmd.Dir = p.RepoDir
	if err := cmd.Run(); err == nil {
		return StatusPending
	}

	return StatusConflict
}
