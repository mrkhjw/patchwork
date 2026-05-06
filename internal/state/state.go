// Package state manages the persistent record of which patches have been
// applied to which repositories, along with their status and timestamps.
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Status values for a patch entry.
const (
	StatusPending = "pending"
	StatusApplied = "applied"
	StatusFailed  = "failed"
	StatusSkipped = "skipped"
)

// Entry records the state of a single patch.
type Entry struct {
	Name      string `json:"name"`
	Repo      string `json:"repo"`
	Status    string `json:"status"`
	AppliedAt string `json:"applied_at,omitempty"`
	Note      string `json:"note,omitempty"`
}

// State holds all patch entries keyed by patch name.
type State struct {
	entries map[string]Entry
}

// New returns an empty State.
func New() *State {
	return &State{entries: make(map[string]Entry)}
}

// Load reads state from a JSON file at path. Returns empty State if missing.
func Load(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("state: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("state: unmarshal: %w", err)
	}
	st := New()
	for _, e := range entries {
		st.entries[e.Name] = e
	}
	return st, nil
}

// Save persists the state to a JSON file at path.
func (s *State) Save(path string) error {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("state: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("state: write: %w", err)
	}
	return nil
}

// Upsert inserts or replaces an entry.
func (s *State) Upsert(e Entry) {
	if e.Status == StatusApplied && e.AppliedAt == "" {
		e.AppliedAt = time.Now().UTC().Format(time.RFC3339)
	}
	s.entries[e.Name] = e
}

// Get returns the entry for patchName and whether it was found.
func (s *State) Get(patchName string) (Entry, bool) {
	e, ok := s.entries[patchName]
	return e, ok
}

// All returns all entries as a slice.
func (s *State) All() []Entry {
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	return list
}
