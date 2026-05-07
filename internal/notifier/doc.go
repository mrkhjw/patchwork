// Package notifier dispatches post-run notifications to one or more
// configured targets after a patchwork run finishes.
//
// Supported target kinds:
//
//   - stdout  — prints the summary to standard output.
//   - file    — appends the summary to the given file path.
//   - exec    — pipes the summary to stdin of the given shell command.
//
// Example configuration (patchwork.yaml):
//
//	notifiers:
//	  - kind: file
//	    destination: /var/log/patchwork.log
//	  - kind: exec
//	    destination: "slack-notify --channel ops"
//
// Notify continues dispatching to remaining targets even when one fails,
// and returns the first error encountered.
package notifier
