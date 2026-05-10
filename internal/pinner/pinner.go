// Package pinner records and enforces pinned patch versions by storing
// a content hash at pin-time and alerting when the patch file changes.
package pinner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Pin holds the recorded hash for a single patch file.
type Pin struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Hash      string    `json:"hash"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// Result is the outcome of verifying a pin.
type Result struct {
	Pin
	Drifted bool   `json:"drifted"`
	Error   string `json:"error,omitempty"`
}

// Store maps patch name to Pin.
type Store map[string]Pin

// PinFile records the current hash of path under name into the store.
func PinFile(store Store, name, path string) error {
	h, err := hashFile(path)
	if err != nil {
		return fmt.Errorf("pinner: hash %s: %w", path, err)
	}
	store[name] = Pin{Name: name, Path: path, Hash: h, PinnedAt: time.Now().UTC()}
	return nil
}

// Verify checks each pin in the store against the current file hash.
func Verify(store Store) []Result {
	results := make([]Result, 0, len(store))
	for _, p := range store {
		r := Result{Pin: p}
		current, err := hashFile(p.Path)
		if err != nil {
			r.Error = err.Error()
			r.Drifted = true
		} else if current != p.Hash {
			r.Drifted = true
		}
		results = append(results, r)
	}
	return results
}

// AnyDrifted returns true if any result has Drifted set.
func AnyDrifted(results []Result) bool {
	for _, r := range results {
		if r.Drifted {
			return true
		}
	}
	return false
}

// SaveStore writes the store to path as JSON.
func SaveStore(store Store, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("pinner: create store: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(store)
}

// LoadStore reads a store from path. Returns an empty store if the file is missing.
func LoadStore(path string) (Store, error) {
	store := make(Store)
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return store, nil
		}
		return nil, fmt.Errorf("pinner: open store: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		return nil, fmt.Errorf("pinner: decode store: %w", err)
	}
	return store, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
