// Package watchdog monitors the health of tracked patches by inspecting
// state entries and flagging those that have remained in a non-terminal
// status (e.g. "pending" or "failed") beyond configurable time thresholds.
//
// # Overview
//
// The primary entry point is [Watch], which accepts a slice of state entries,
// a reference time (typically time.Now()), and a [Policy] describing warn/crit
// age thresholds and which statuses to monitor.
//
// # Alerts
//
// Each stale entry produces an [Alert] with a [Severity] of either "warn" or
// "critical". Use [AnyFailed] to check whether any critical alerts exist, and
// [Format] to render a human-readable table suitable for CLI output.
//
// # Policy
//
// [DefaultPolicy] watches "pending" and "failed" statuses, warning after 24 h
// and escalating to critical after 72 h.
package watchdog
