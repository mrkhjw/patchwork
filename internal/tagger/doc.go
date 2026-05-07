// Package tagger builds and queries an inverted index of tags to patch names.
//
// Given a list of patch definitions, tagger.Build returns a TagIndex that
// lets callers quickly answer questions such as:
//
//	"Which patches are tagged 'security'?"
//	"What tags does patch 'fix-auth' carry?"
//
// The index is read-only after construction and safe for concurrent reads.
// All returned slices are sorted lexicographically for stable output.
//
// Typical usage:
//
//	cfg, _ := config.Load("patchwork.yaml")
//	idx  := tagger.Build(cfg.Patches)
//	for _, tag := range idx.AllTags() {
//		fmt.Println(tag, idx.PatchesForTag(tag))
//	}
package tagger
