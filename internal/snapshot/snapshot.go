// Package snapshot captures and compares the state of patches across repos
// at a point in time, enabling drift detection and audit trails.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents the state of a single patch in a single repo at snapshot time.
type Entry struct {
	PatchName string    `json:"patch_name"`
	Repo      string    `json:"repo"`
	Applied   bool      `json:"applied"`
	Hash      string    `json:"hash,omitempty"`
	CapturedAt time.Time `json:"captured_at"`
}

// Snapshot holds a collection of entries captured at a point in time.
type Snapshot struct {
	CapturedAt time.Time `json:"captured_at"`
	Entries    []Entry   `json:"entries"`
}

// New creates a new Snapshot with the given entries, stamping the current time.
func New(entries []Entry) *Snapshot {
	now := time.Now().UTC()
	for i := range entries {
		entries[i].CapturedAt = now
	}
	return &Snapshot{
		CapturedAt: now,
		Entries:    entries,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(path string, s *Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Diff returns entries present in current but whose Applied status differs from baseline.
func Diff(baseline, current *Snapshot) []Entry {
	index := make(map[string]bool, len(baseline.Entries))
	for _, e := range baseline.Entries {
		key := e.Repo + "|" + e.PatchName
		index[key] = e.Applied
	}
	var drifted []Entry
	for _, e := range current.Entries {
		key := e.Repo + "|" + e.PatchName
		if applied, found := index[key]; !found || applied != e.Applied {
			drifted = append(drifted, e)
		}
	}
	return drifted
}
