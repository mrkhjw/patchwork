// Package exporter provides utilities for exporting patch run results to
// external formats such as JSON and CSV.
//
// Usage:
//
//	entries := []exporter.Entry{
//		{PatchName: "fix-config", Repo: "api", Status: "applied", AppliedAt: time.Now()},
//	}
//
//	// Write JSON to stdout
//	exporter.Export(os.Stdout, entries, exporter.FormatJSON)
//
//	// Write CSV to a file
//	f, _ := os.Create("results.csv")
//	defer f.Close()
//	exporter.Export(f, entries, exporter.FormatCSV)
//
// Supported formats: "json", "csv".
package exporter
