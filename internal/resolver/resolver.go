// Package resolver resolves patch file paths relative to a base directory,
// supporting both absolute and relative path specifications.
package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

// Result holds the resolved absolute path and whether the file exists.
type Result struct {
	PatchName    string
	OriginalPath string
	ResolvedPath string
	Exists       bool
}

// Resolve takes a base directory and a list of (name, path) pairs and returns
// resolved Results. Relative paths are joined with baseDir.
func Resolve(baseDir string, patches map[string]string) ([]Result, error) {
	if baseDir == "" {
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("resolver: cannot determine working directory: %w", err)
		}
	}

	results := make([]Result, 0, len(patches))
	for name, p := range patches {
		resolved := p
		if !filepath.IsAbs(p) {
			resolved = filepath.Join(baseDir, p)
		}
		resolved = filepath.Clean(resolved)

		_, err := os.Stat(resolved)
		exists := err == nil

		results = append(results, Result{
			PatchName:    name,
			OriginalPath: p,
			ResolvedPath: resolved,
			Exists:       exists,
		})
	}
	return results, nil
}

// AnyMissing returns true if at least one resolved result does not exist on disk.
func AnyMissing(results []Result) bool {
	for _, r := range results {
		if !r.Exists {
			return true
		}
	}
	return false
}

// Missing returns only the results where the file does not exist.
func Missing(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if !r.Exists {
			out = append(out, r)
		}
	}
	return out
}
