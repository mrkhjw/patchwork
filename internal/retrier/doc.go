// Package retrier provides configurable retry logic for patchwork operations.
//
// A Policy defines the maximum number of attempts, the initial delay between
// retries, and an optional exponential backoff multiplier.
//
// Usage:
//
//	p := retrier.DefaultPolicy()
//	r := retrier.Run(p, func() error {
//	    return applyPatch()
//	}, nil)
//	if r.Err != nil {
//	    log.Printf("patch failed after %d attempts: %v", r.Attempts, r.Err)
//	}
//
// For errors that should not be retried (e.g. patch file not found), wrap
// them with retrier.Permanent so that RunSmart stops immediately:
//
//	return retrier.Permanent{Err: fmt.Errorf("patch file missing: %s", path)}
package retrier
