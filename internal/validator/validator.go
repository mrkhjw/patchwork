// Package validator provides pre-flight checks for patches and repositories
// before the runner attempts to apply them.
package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/repo"
)

// Result holds the outcome of validating a single patch entry.
type Result struct {
	PatchName string
	RepoPath  string
	Errors    []string
}

// OK returns true when no validation errors were found.
func (r Result) OK() bool {
	return len(r.Errors) == 0
}

// Validate runs pre-flight checks on all patches defined in cfg.
// It verifies that each referenced repository exists, is on the expected
// branch (when specified), and that the patch file itself is readable.
func Validate(cfg *config.Config) []Result {
	var results []Result

	for _, p := range cfg.Patches {
		for _, repoPath := range p.Repos {
			r := Result{
				PatchName: p.Name,
				RepoPath:  repoPath,
			}

			if !repo.Exists(repoPath) {
				r.Errors = append(r.Errors, fmt.Sprintf("repository not found: %s", repoPath))
			} else if branch, err := repo.CurrentBranch(repoPath); err != nil {
				r.Errors = append(r.Errors, fmt.Sprintf("could not determine branch: %v", err))
			} else if p.Branch != "" && branch != p.Branch {
				r.Errors = append(r.Errors,
					fmt.Sprintf("expected branch %q but found %q", p.Branch, branch))
			}

			patchFile := filepath.Clean(p.File)
			if _, err := os.Stat(patchFile); err != nil {
				r.Errors = append(r.Errors, fmt.Sprintf("patch file unreadable: %s", patchFile))
			}

			results = append(results, r)
		}
	}

	return results
}

// AnyFailed returns true if at least one Result contains errors.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if !r.OK() {
			return true
		}
	}
	return false
}
