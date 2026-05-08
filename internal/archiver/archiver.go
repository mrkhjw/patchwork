// Package archiver compresses and stores patch-related artifacts (logs, diffs,
// audit trails) into a single zip archive for long-term retention.
package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Result holds the outcome of a single file archival operation.
type Result struct {
	Source  string
	Archive string
	Err     error
}

// Archive compresses the given source files into a zip archive written to
// destDir. The archive is named using the provided label and a UTC timestamp.
// It returns one Result per source file.
func Archive(destDir, label string, sources []string) ([]Result, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("archiver: create dest dir: %w", err)
	}

	stamp := time.Now().UTC().Format("20060102T150405Z")
	archiveName := fmt.Sprintf("%s-%s.zip", label, stamp)
	archivePath := filepath.Join(destDir, archiveName)

	zf, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("archiver: create archive: %w", err)
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	results := make([]Result, 0, len(sources))
	for _, src := range sources {
		r := Result{Source: src, Archive: archivePath}
		if err := addFile(zw, src); err != nil {
			r.Err = err
		}
		results = append(results, r)
	}
	return results, nil
}

// AnyFailed returns true if any Result contains a non-nil error.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

func addFile(zw *zip.Writer, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer f.Close()

	w, err := zw.Create(filepath.Base(src))
	if err != nil {
		return fmt.Errorf("zip entry %s: %w", src, err)
	}

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("copy %s: %w", src, err)
	}
	return nil
}
