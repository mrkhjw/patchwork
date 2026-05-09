// Package labeler derives labels for patches automatically from their metadata.
//
// Unlike tags (which are author-declared in the config), labels are computed at
// runtime by evaluating a list of Rule values against each PatchMeta. Rules can
// match on repo path substrings, patch status, or how long ago the patch was
// applied.
//
// Typical usage:
//
//	rules := []labeler.Rule{
//		{Label: "stale",   OlderThanDays: 30},
//		{Label: "pending", StatusEquals:  "pending"},
//		{Label: "infra",   RepoContains:  "infra"},
//	}
//	results := labeler.Label(patches, rules)
//	for _, r := range results {
//		fmt.Println(labeler.Format(r))
//	}
package labeler
