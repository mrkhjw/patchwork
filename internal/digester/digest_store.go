package digester

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Store is a persisted map of patch name → SHA-256 digest.
type Store map[string]string

// SaveStore writes the digest store to path as JSON.
func SaveStore(path string, s Store) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("digester: marshal store: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("digester: write store: %w", err)
	}
	return nil
}

// LoadStore reads a digest store from path.
// Returns an empty store when the file does not exist.
func LoadStore(path string) (Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Store{}, nil
		}
		return nil, fmt.Errorf("digester: read store: %w", err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("digester: parse store: %w", err)
	}
	return s, nil
}

// StoreFromResults builds a Store from a slice of successful Results.
func StoreFromResults(results []Result) Store {
	s := make(Store, len(results))
	for _, r := range results {
		if r.OK() {
			s[r.PatchName] = r.Digest
		}
	}
	return s
}

// Changed returns patch names whose digest differs from the stored value.
// Patches absent from the store are not reported as changed.
func Changed(current Store, previous Store) []string {
	var drifted []string
	for name, digest := range current {
		if prev, ok := previous[name]; ok && prev != digest {
			drifted = append(drifted, name)
		}
	}
	return drifted
}
