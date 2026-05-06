// Package audit records a timestamped log of every patch operation
// (apply, skip, rollback) so operators can trace what happened and when.
package audit

import (
	"encoding/json"
	"os"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventApplied  EventKind = "applied"
	EventSkipped  EventKind = "skipped"
	EventRolledBack EventKind = "rolled_back"
	EventFailed   EventKind = "failed"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Patch     string    `json:"patch"`
	Repo      string    `json:"repo"`
	Kind      EventKind `json:"kind"`
	Message   string    `json:"message,omitempty"`
}

// Log holds an ordered list of audit events.
type Log struct {
	Events []Event `json:"events"`
}

// Append adds a new event to the log, stamping the current UTC time.
func (l *Log) Append(patch, repo string, kind EventKind, message string) {
	l.Events = append(l.Events, Event{
		Timestamp: time.Now().UTC(),
		Patch:     patch,
		Repo:      repo,
		Kind:      kind,
		Message:   message,
	})
}

// Save writes the log as JSON to path, creating or truncating the file.
func (l *Log) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(l)
}

// Load reads an audit log from path. Returns an empty Log if the file does
// not exist.
func Load(path string) (*Log, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Log{}, nil
	}
	if err != nil {
		return nil, err
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Filter returns only the events that match the given kind. Pass an empty
// string to return all events.
func (l *Log) Filter(kind EventKind) []Event {
	if kind == "" {
		return l.Events
	}
	var out []Event
	for _, e := range l.Events {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}
