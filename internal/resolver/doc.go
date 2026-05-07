// Package resolver provides utilities for resolving patch file paths used by
// patchwork configurations.
//
// Patches may be specified with absolute or relative paths in the YAML config.
// The resolver normalises these against a configurable base directory (defaulting
// to the current working directory) and reports which files are present on disk.
//
// Typical usage:
//
//	results, err := resolver.Resolve("/home/user/patches", map[string]string{
//		"fix-auth": "fix-auth.patch",
//		"bump-deps": "/absolute/path/bump.patch",
//	})
//	if resolver.AnyMissing(results) {
//		for _, r := range resolver.Missing(results) {
//			fmt.Println("missing:", r.ResolvedPath)
//		}
//	}
package resolver
