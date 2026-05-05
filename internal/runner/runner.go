package runner

import (
	"fmt"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/patch"
	"github.com/patchwork/internal/repo"
	"github.com/patchwork/internal/state"
)

// Result holds the outcome of applying a patch to a single repo.
type Result struct {
	Repo    string
	Patch   string
	Applied bool
	Skipped bool
	Err     error
}

// RunOptions controls runner behaviour.
type RunOptions struct {
	DryRun   bool
	Force    bool
	StatePath string
}

// Run iterates over all repos defined in cfg and applies any pending patches,
// recording results in the state file located at opts.StatePath.
func Run(cfg *config.Config, opts RunOptions) ([]Result, error) {
	st, err := state.Load(opts.StatePath)
	if err != nil {
		return nil, fmt.Errorf("runner: load state: %w", err)
	}

	var results []Result

	for _, r := range cfg.Repos {
		if !repo.Exists(r.Path) {
			results = append(results, Result{
				Repo:  r.Name,
				Err:   fmt.Errorf("repo path does not exist: %s", r.Path),
			})
			continue
		}

		for _, p := range cfg.Patches {
			entry, _ := st.Get(r.Name, p.Name)
			if entry != nil && entry.Status == "applied" && !opts.Force {
				results = append(results, Result{
					Repo:    r.Name,
					Patch:   p.Name,
					Skipped: true,
				})
				continue
			}

			if opts.DryRun {
				pending, err := patch.Check(r.Path, p.File)
				results = append(results, Result{
					Repo:    r.Name,
					Patch:   p.Name,
					Applied: !pending,
					Skipped: !pending,
					Err:     err,
				})
				continue
			}

			applyErr := patch.Apply(r.Path, p.File)
			status := "applied"
			if applyErr != nil {
				status = "failed"
			}
			st.Upsert(r.Name, p.Name, status)
			results = append(results, Result{
				Repo:    r.Name,
				Patch:   p.Name,
				Applied: applyErr == nil,
				Err:     applyErr,
			})
		}
	}

	if !opts.DryRun {
		if err := st.Save(opts.StatePath); err != nil {
			return results, fmt.Errorf("runner: save state: %w", err)
		}
	}

	return results, nil
}
