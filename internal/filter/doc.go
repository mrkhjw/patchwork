// Package filter provides patch selection logic for patchwork.
//
// Callers build an Options struct describing which patches they care about
// (by name, tag, or target repo) and pass it together with a full patch list
// to Apply. Apply returns only the patches that satisfy every non-empty
// criterion, making it easy to run a targeted subset of patches without
// modifying the underlying configuration.
//
// Example:
//
//	matched := filter.Apply(cfg.Patches, filter.Options{
//		Tags:  []string{"hotfix"},
//		Repos: []string{"api"},
//	})
package filter
