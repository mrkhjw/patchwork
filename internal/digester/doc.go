// Package digester provides SHA-256 content hashing for patch files.
//
// It is used by the patchwork pipeline to detect file tampering or
// accidental modification before a patch is applied to a repository.
//
// Typical usage:
//
//	results := digester.Digest(cfg.Patches)
//	if digester.AnyFailed(results) {
//		// handle errors
//	}
//	for _, r := range results {
//		fmt.Println(r.PatchName, r.Digest)
//	}
//
// Verify can be used to assert that a previously recorded digest still
// matches the current file content:
//
//	if err := digester.Verify(path, savedDigest); err != nil {
//		// file has changed
//	}
package digester
