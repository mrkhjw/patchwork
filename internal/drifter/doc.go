// Package drifter compares a previously captured snapshot of patch state
// against the current live state to identify unexpected changes — known as
// "drift" — that may indicate manual intervention, failed rollbacks, or
// out-of-band modifications to a repository.
//
// Typical usage:
//
//	snap, _ := snapshot.Load(snapshotPath)
//	st, _   := state.Load(statePath)
//	result  := drifter.Detect(snap, st)
//	if result.AnyDrifted() {
//	    fmt.Println(drifter.Format(result))
//	}
//
// Alerts are produced for missing entries, status changes, and repo path
// changes. Each alert records the patch name, affected field, expected and
// actual values, and the UTC timestamp of detection.
package drifter
