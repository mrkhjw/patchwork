// Package filter provides patch selection logic for patchwork.
//
// Callers build an Options struct describing which patches they care about
// (by name, tag, or target repo) and pass it together with a full patch list
// to Apply. Apply returns only the patches that satisfy every non-empty
// criterion, making it easy to run a targeted subset of patches without
// modifying the underlying configuration.
//
// Filtering rules:
//   - Names: patch name must match one of the provided names (exact match).
//   - Tags: patch must carry at least one of the provided tags.
//   - Repos: patch must target at least one of the provided repos.
//
// When multiple criteria are set, a patch must satisfy all of them (AND
// semantics). An empty Options struct passes every patch through unchanged.
//
// Example:
//
//	matched := filter.Apply(cfg.Patches, filter.Options{
//		Tags:  []string{"hotfix"},
//		Repos: []string{"api"},
//	})
package filter
