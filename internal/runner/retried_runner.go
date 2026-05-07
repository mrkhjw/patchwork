package runner

import (
	"fmt"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/retrier"
)

// RetriedOptions extends Options with retry configuration.
type RetriedOptions struct {
	Options
	Retry retrier.Policy
}

// RunWithRetry behaves like Run but wraps each patch application in retry
// logic governed by opts.Retry. Permanent errors skip remaining attempts.
func RunWithRetry(cfg *config.Config, opts RetriedOptions) []Result {
	results := make([]Result, 0, len(cfg.Patches))

	for i := range cfg.Patches {
		patch := cfg.Patches[i]

		// Run the single-patch slice through the standard runner with retries.
		var last []Result
		retryResult := retrier.RunSmart(opts.Retry, func() error {
			slice := &config.Config{
				Repos:   cfg.Repos,
				Patches: []config.Patch{patch},
			}
			last = Run(slice, opts.Options)
			if len(last) == 0 {
				return nil
			}
			if last[0].Err != nil {
				return last[0].Err
			}
			return nil
		}, nil)

		if len(last) > 0 {
			if retryResult.Err != nil && last[0].Err == nil {
				// Overwrite with the final retry error for clarity.
				last[0].Err = fmt.Errorf("retry exhausted: %w", retryResult.Err)
			}
			results = append(results, last...)
		}
	}

	return results
}
