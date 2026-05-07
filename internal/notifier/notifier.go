// Package notifier provides simple notification hooks that fire after
// a patchwork run completes, summarising results to configured targets.
package notifier

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/yourorg/patchwork/internal/reporter"
)

// Target describes a single notification destination.
type Target struct {
	// Kind is the notification type: "stdout", "file", or "exec".
	Kind string `yaml:"kind"`
	// Destination is the file path (for "file") or command (for "exec").
	Destination string `yaml:"destination,omitempty"`
}

// Notify sends the run summary to every configured target.
// It returns the first error encountered but continues notifying remaining targets.
func Notify(targets []Target, entries []reporter.Entry) error {
	summary := buildSummary(entries)
	var firstErr error
	for _, t := range targets {
		if err := dispatch(t, summary); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("notifier: target %q: %w", t.Kind, err)
		}
	}
	return firstErr
}

// AnyFailed returns true when at least one entry has a failed status.
func AnyFailed(entries []reporter.Entry) bool {
	for _, e := range entries {
		if e.Status == "failed" {
			return true
		}
	}
	return false
}

func dispatch(t Target, summary string) error {
	switch t.Kind {
	case "stdout":
		_, err := fmt.Fprint(os.Stdout, summary)
		return err
	case "file":
		return writeFile(t.Destination, summary)
	case "exec":
		return runCommand(t.Destination, summary)
	default:
		return fmt.Errorf("unknown kind %q", t.Kind)
	}
}

func writeFile(path, summary string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, summary)
	return err
}

func runCommand(command, summary string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec
	cmd.Stdin = strings.NewReader(summary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildSummary(entries []reporter.Entry) string {
	s := reporter.Summary(entries)
	return fmt.Sprintf("patchwork run complete — applied: %d, skipped: %d, failed: %d\n",
		s.Applied, s.Skipped, s.Failed)
}
