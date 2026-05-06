// Package hooks provides pre- and post-apply hook execution for patches.
package hooks

import (
	"fmt"
	"os/exec"
	"strings"
)

// Hook represents a shell command to run at a specific lifecycle point.
type Hook struct {
	Command string
	Dir     string
}

// RunResult captures the outcome of a single hook execution.
type RunResult struct {
	Command string
	Output  string
	Err     error
}

// RunAll executes each hook in order, stopping on the first failure.
// dir is used as the working directory when the hook has no Dir set.
func RunAll(hooks []Hook, defaultDir string) []RunResult {
	results := make([]RunResult, 0, len(hooks))
	for _, h := range hooks {
		result := run(h, defaultDir)
		results = append(results, result)
		if result.Err != nil {
			break
		}
	}
	return results
}

// AnyFailed returns true if any result contains a non-nil error.
func AnyFailed(results []RunResult) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

func run(h Hook, defaultDir string) RunResult {
	dir := h.Dir
	if dir == "" {
		dir = defaultDir
	}
	parts := strings.Fields(h.Command)
	if len(parts) == 0 {
		return RunResult{Command: h.Command, Err: fmt.Errorf("empty command")}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return RunResult{
		Command: h.Command,
		Output:  strings.TrimSpace(string(out)),
		Err:     err,
	}
}
