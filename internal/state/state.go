package state

import (
	"encoding/json"
	"os"
	"time"
)

// PatchStatus represents the application status of a single patch on a repo.
type PatchStatus struct {
	PatchName  string    `json:"patch_name"`
	Repo       string    `json:"repo"`
	Applied    bool      `json:"applied"`
	AppliedAt  time.Time `json:"applied_at,omitempty"`
	Checksum   string    `json:"checksum"`
}

// State holds the full tracked state for all patches across all repos.
type State struct {
	Entries []PatchStatus `json:"entries"`
}

// Load reads state from the given JSON file path.
// Returns an empty State if the file does not exist.
func Load(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save writes the state to the given JSON file path.
func (s *State) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Upsert adds or updates the status entry for a given patch+repo pair.
func (s *State) Upsert(ps PatchStatus) {
	for i, e := range s.Entries {
		if e.PatchName == ps.PatchName && e.Repo == ps.Repo {
			s.Entries[i] = ps
			return
		}
	}
	s.Entries = append(s.Entries, ps)
}

// Get returns the PatchStatus for a given patch+repo pair, and whether it was found.
func (s *State) Get(patchName, repo string) (PatchStatus, bool) {
	for _, e := range s.Entries {
		if e.PatchName == patchName && e.Repo == repo {
			return e, true
		}
	}
	return PatchStatus{}, false
}
