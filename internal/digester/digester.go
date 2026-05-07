// Package digester computes and verifies content hashes for patch files,
// enabling tamper detection and integrity checks before application.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/patchwork/internal/config"
)

// Result holds the computed digest for a single patch.
type Result struct {
	PatchName string
	FilePath  string
	Digest    string
	Err       error
}

// OK returns true when no error occurred.
func (r Result) OK() bool { return r.Err == nil }

// Digest computes SHA-256 hashes for all patch files referenced in patches.
func Digest(patches []config.Patch) []Result {
	results := make([]Result, 0, len(patches))
	for _, p := range patches {
		results = append(results, digestOne(p))
	}
	return results
}

// Verify checks that the SHA-256 of the file at path matches expected.
func Verify(path, expected string) error {
	got, err := hashFile(path)
	if err != nil {
		return fmt.Errorf("digester: cannot read %q: %w", path, err)
	}
	if got != expected {
		return fmt.Errorf("digester: digest mismatch for %q: want %s got %s", path, expected, got)
	}
	return nil
}

// AnyFailed returns true if any result contains an error.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if !r.OK() {
			return true
		}
	}
	return false
}

func digestOne(p config.Patch) Result {
	h, err := hashFile(p.Path)
	if err != nil {
		return Result{PatchName: p.Name, FilePath: p.Path, Err: fmt.Errorf("digester: %w", err)}
	}
	return Result{PatchName: p.Name, FilePath: p.Path, Digest: h}
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
