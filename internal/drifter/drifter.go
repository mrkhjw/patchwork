// Package drifter detects configuration drift between the expected patch
// state recorded in a snapshot and the current live state of each repo.
package drifter

import (
	"fmt"
	"time"

	"github.com/patchwork/internal/snapshot"
	"github.com/patchwork/internal/state"
)

// Alert describes a single drift event for one patch entry.
type Alert struct {
	PatchName string
	Repo      string
	Field     string
	Expected  string
	Actual    string
	DetectedAt time.Time
}

// Result holds all alerts produced by a drift check.
type Result struct {
	Alerts []Alert
}

// AnyDrifted returns true when at least one alert was produced.
func (r Result) AnyDrifted() bool { return len(r.Alerts) > 0 }

// Detect compares a previously saved snapshot against the current state
// and returns a Result describing every field that has changed.
func Detect(snap snapshot.Snapshot, st *state.State) Result {
	var alerts []Alert
	now := time.Now().UTC()

	for _, entry := range snap.Entries {
		current, ok := st.Get(entry.PatchName)
		if !ok {
			alerts = append(alerts, Alert{
				PatchName:  entry.PatchName,
				Repo:       entry.Repo,
				Field:      "existence",
				Expected:   "present",
				Actual:     "missing",
				DetectedAt: now,
			})
			continue
		}
		if current.Status != entry.Status {
			alerts = append(alerts, Alert{
				PatchName:  entry.PatchName,
				Repo:       entry.Repo,
				Field:      "status",
				Expected:   entry.Status,
				Actual:     current.Status,
				DetectedAt: now,
			})
		}
		if current.Repo != entry.Repo {
			alerts = append(alerts, Alert{
				PatchName:  entry.PatchName,
				Repo:       entry.Repo,
				Field:      "repo",
				Expected:   entry.Repo,
				Actual:     current.Repo,
				DetectedAt: now,
			})
		}
	}
	return Result{Alerts: alerts}
}

// Format returns a human-readable summary of the drift result.
func Format(r Result) string {
	if !r.AnyDrifted() {
		return "no drift detected"
	}
	out := fmt.Sprintf("%d drift alert(s):\n", len(r.Alerts))
	for _, a := range r.Alerts {
		out += fmt.Sprintf("  [%s] %s.%s: expected=%q actual=%q\n",
			a.Repo, a.PatchName, a.Field, a.Expected, a.Actual)
	}
	return out
}
