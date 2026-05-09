// Package watchdog monitors patch health over time, flagging entries
// that have been pending or failed beyond a configurable threshold.
package watchdog

import (
	"fmt"
	"time"

	"github.com/patchwork/internal/state"
)

// Severity indicates how critical a watchdog alert is.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityCrit  Severity = "critical"
)

// Alert represents a single watchdog finding for a patch entry.
type Alert struct {
	PatchName string
	Repo      string
	Status    string
	Age       time.Duration
	Severity  Severity
	Message   string
}

// Policy controls thresholds for alert generation.
type Policy struct {
	WarnAfter time.Duration
	CritAfter time.Duration
	WatchStatuses []string
}

// DefaultPolicy returns a sensible default watchdog policy.
func DefaultPolicy() Policy {
	return Policy{
		WarnAfter:     24 * time.Hour,
		CritAfter:     72 * time.Hour,
		WatchStatuses: []string{"pending", "failed"},
	}
}

// Watch inspects state entries and returns alerts for stale patches.
func Watch(entries []state.Entry, now time.Time, p Policy) []Alert {
	watched := toSet(p.WatchStatuses)
	var alerts []Alert

	for _, e := range entries {
		if !watched[e.Status] {
			continue
		}
		age := now.Sub(e.UpdatedAt)
		var sev Severity
		switch {
		case age >= p.CritAfter:
			sev = SeverityCrit
		case age >= p.WarnAfter:
			sev = SeverityWarn
		default:
			continue
		}
		alerts = append(alerts, Alert{
			PatchName: e.PatchName,
			Repo:      e.Repo,
			Status:    e.Status,
			Age:       age.Truncate(time.Second),
			Severity:  sev,
			Message:   fmt.Sprintf("patch %q in repo %q has been %s for %s", e.PatchName, e.Repo, e.Status, age.Truncate(time.Second)),
		})
	}
	return alerts
}

// AnyFailed returns true if any alert is critical severity.
func AnyFailed(alerts []Alert) bool {
	for _, a := range alerts {
		if a.Severity == SeverityCrit {
			return true
		}
	}
	return false
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}
