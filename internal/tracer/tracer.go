// Package tracer records per-patch execution timing and outcome for
// post-run diagnostics and performance analysis.
package tracer

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Span holds the timing and result of a single patch execution.
type Span struct {
	PatchName string
	Repo      string
	Started   time.Time
	Duration  time.Duration
	Err       error
}

// Tracer collects Spans during a run.
type Tracer struct {
	spans []Span
}

// New returns an initialised Tracer.
func New() *Tracer {
	return &Tracer{}
}

// Start records the current time and returns a function that, when called,
// appends a completed Span to the Tracer.
func (t *Tracer) Start(patchName, repo string) func(err error) {
	begin := time.Now()
	return func(err error) {
		t.spans = append(t.spans, Span{
			PatchName: patchName,
			Repo:      repo,
			Started:   begin,
			Duration:  time.Since(begin),
			Err:       err,
		})
	}
}

// Spans returns a copy of all recorded spans.
func (t *Tracer) Spans() []Span {
	out := make([]Span, len(t.spans))
	copy(out, t.spans)
	return out
}

// Slowest returns the n spans with the longest duration, descending.
func (t *Tracer) Slowest(n int) []Span {
	sorted := t.Spans()
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration > sorted[j].Duration
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// AnyFailed reports whether at least one span recorded a non-nil error.
func (t *Tracer) AnyFailed() bool {
	for _, s := range t.spans {
		if s.Err != nil {
			return true
		}
	}
	return false
}

// Write formats all spans as a human-readable table into w.
func (t *Tracer) Write(w io.Writer) {
	fmt.Fprintf(w, "%-30s %-25s %10s  %s\n", "PATCH", "REPO", "DURATION", "STATUS")
	for _, s := range t.spans {
		status := "ok"
		if s.Err != nil {
			status = "FAILED: " + s.Err.Error()
		}
		fmt.Fprintf(w, "%-30s %-25s %10s  %s\n",
			s.PatchName, s.Repo, s.Duration.Round(time.Millisecond), status)
	}
}
