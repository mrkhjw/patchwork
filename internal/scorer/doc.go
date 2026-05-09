// Package scorer computes a numeric priority score for each patch defined in
// the patchwork configuration.
//
// Scores are derived from three sources:
//
//  1. Tag weights – configurable bonuses assigned to specific tags such as
//     "critical" or "hotfix". A patch inherits the sum of weights for all
//     matching tags.
//
//  2. Age bonus – patches that were last applied longer ago receive a small
//     additive bonus proportional to the number of elapsed days, encouraging
//     periodic re-evaluation of long-standing patches.
//
//  3. Retry penalty – each failed application attempt subtracts a fixed
//     amount from the score, surfacing consistently problematic patches for
//     manual review.
//
// Callers can supply a custom Policy to override any of these weights, or use
// DefaultPolicy for sensible out-of-the-box behaviour.
//
// The AnyNegative helper makes it easy to detect patches whose accumulated
// penalties have driven their score below zero.
package scorer
