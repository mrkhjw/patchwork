// Package hooks provides lifecycle hook execution for patchwork operations.
//
// Hooks are shell commands associated with a patch or global configuration
// that run before or after a patch is applied. They can be used to:
//
//   - Validate environment prerequisites
//   - Restart services after a config patch is applied
//   - Notify external systems of changes
//
// Hooks are executed in order; if any hook exits with a non-zero status,
// subsequent hooks in the same list are skipped and the error is surfaced
// to the caller.
//
// Example usage:
//
//	hooks.RunAll([]hooks.Hook{
//		{Command: "make lint"},
//		{Command: "systemctl reload myservice"},
//	}, repoDir)
package hooks
