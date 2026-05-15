// Package classifier assigns risk tiers (critical, high, medium, low) to
// patches based on their tags and current state.
//
// Usage:
//
//	pol := classifier.DefaultPolicy()
//	results := classifier.Classify(cfg, st, pol)
//	if classifier.AnyAbove(results, classifier.TierHigh) {
//		fmt.Println("critical patches detected")
//	}
//	fmt.Print(classifier.Format(results))
//
// The default policy flags patches tagged "security" or "breaking" as
// critical, "migration" or "schema" as high, and any patch whose recorded
// status is "failed" as at least high risk.
package classifier
