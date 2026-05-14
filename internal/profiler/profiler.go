// Package profiler measures execution time and resource usage for patch operations.
package profiler

import (
	"fmt"
	"sort"
	"time"
)

// Record holds timing data for a single named operation.
type Record struct {
	Name     string
	Duration time.Duration
	Err      error
}

// Profile holds all records collected during a run.
type Profile struct {
	Records []Record
}

// New returns an empty Profile.
func New() *Profile {
	return &Profile{}
}

// Track runs fn, records its duration under name, and stores any error.
func (p *Profile) Track(name string, fn func() error) error {
	start := time.Now()
	err := fn()
	p.Records = append(p.Records, Record{
		Name:     name,
		Duration: time.Since(start),
		Err:      err,
	})
	return err
}

// Slowest returns up to n records sorted by duration descending.
func (p *Profile) Slowest(n int) []Record {
	sorted := make([]Record, len(p.Records))
	copy(sorted, p.Records)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration > sorted[j].Duration
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// AnyFailed reports whether any tracked operation returned a non-nil error.
func (p *Profile) AnyFailed() bool {
	for _, r := range p.Records {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// TotalDuration returns the sum of all recorded durations.
func (p *Profile) TotalDuration() time.Duration {
	var total time.Duration
	for _, r := range p.Records {
		total += r.Duration
	}
	return total
}

// Format returns a human-readable summary of the profile.
func Format(p *Profile) string {
	if len(p.Records) == 0 {
		return "profiler: no records\n"
	}
	out := fmt.Sprintf("profiler: %d operations, total %s\n", len(p.Records), p.TotalDuration().Round(time.Millisecond))
	for _, r := range p.Records {
		status := "ok"
		if r.Err != nil {
			status = fmt.Sprintf("err: %s", r.Err)
		}
		out += fmt.Sprintf("  %-40s %8s  [%s]\n", r.Name, r.Duration.Round(time.Millisecond), status)
	}
	return out
}
