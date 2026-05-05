// Package runner orchestrates the application of patches across multiple
// repositories as described by a patchwork configuration.
//
// It ties together the config, patch, repo, and state packages:
//
//  1. For each repo declared in the config, runner verifies the repo exists
//     on disk using the repo package.
//
//  2. For each patch declared in the config, runner consults the state package
//     to determine whether the patch has already been applied to that repo.
//
//  3. If the patch is pending (or --force is set), runner delegates to the
//     patch package to apply the diff and then persists the new status via
//     the state package.
//
// Dry-run mode performs all checks but skips both the apply step and any
// state mutations, making it safe to use in CI preview workflows.
package runner
