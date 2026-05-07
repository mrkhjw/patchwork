// Package linter validates patch configuration entries for common
// mistakes such as empty fields, invalid file paths, and unknown repos.
package linter

import (
	"fmt"
	"os"
	"strings"

	"github.com/patchwork/internal/config"
)

// Result holds the outcome of linting a single patch entry.
type Result struct {
	PatchName string
	Errors    []string
	Warnings  []string
}

// AnyErrors returns true if the result contains at least one error.
func (r Result) AnyErrors() bool { return len(r.Errors) > 0 }

// AnyWarnings returns true if the result contains at least one warning.
func (r Result) AnyWarnings() bool { return len(r.Warnings) > 0 }

// Lint checks every patch in cfg and returns one Result per patch.
func Lint(cfg *config.Config) []Result {
	results := make([]Result, 0, len(cfg.Patches))
	for _, p := range cfg.Patches {
		results = append(results, lintPatch(p))
	}
	return results
}

// AnyFailed returns true when at least one result carries errors.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if r.AnyErrors() {
			return true
		}
	}
	return false
}

func lintPatch(p config.Patch) Result {
	r := Result{PatchName: p.Name}

	if strings.TrimSpace(p.Name) == "" {
		r.Errors = append(r.Errors, "patch name must not be empty")
	}

	if strings.TrimSpace(p.File) == "" {
		r.Errors = append(r.Errors, "patch file path must not be empty")
	} else if _, err := os.Stat(p.File); os.IsNotExist(err) {
		r.Errors = append(r.Errors, fmt.Sprintf("patch file not found: %s", p.File))
	}

	if strings.TrimSpace(p.Repo) == "" {
		r.Errors = append(r.Errors, "repo path must not be empty")
	}

	if len(p.Tags) == 0 {
		r.Warnings = append(r.Warnings, "no tags defined; patch will not be reachable by tag filters")
	}

	return r
}
