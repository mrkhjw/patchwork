// Package planner determines which patches should be applied to which repos,
// taking into account current state, filters, and validation results.
package planner

import (
	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/filter"
	"github.com/yourorg/patchwork/internal/state"
	"github.com/yourorg/patchwork/internal/validator"
)

// Task represents a single patch-to-repo application unit.
type Task struct {
	Patch  config.Patch
	Repo   string
	Reason string // why this task was included
}

// Options controls planner behaviour.
type Options struct {
	FilterNames []string
	FilterTags  []string
	FilterRepos []string
	SkipApplied bool
}

// Plan builds an ordered list of Tasks from the provided config, state, and
// validation results, applying any active filters.
func Plan(
	patches []config.Patch,
	st *state.State,
	results []validator.Result,
	opts Options,
) []Task {
	valid := validSet(results)

	filtered := filter.Apply(patches, filter.Options{
		Names: opts.FilterNames,
		Tags:  opts.FilterTags,
		Repos: opts.FilterRepos,
	})

	var tasks []Task
	for _, p := range filtered {
		if !valid[p.Name] {
			continue
		}
		for _, repo := range p.Repos {
			if opts.SkipApplied {
				if entry, ok := st.Get(p.Name, repo); ok && entry.Status == "applied" {
					continue
				}
			}
			reason := "pending"
			if entry, ok := st.Get(p.Name, repo); ok {
				reason = entry.Status
			}
			tasks = append(tasks, Task{
				Patch:  p,
				Repo:   repo,
				Reason: reason,
			})
		}
	}
	return tasks
}

// validSet returns a set of patch names that passed validation.
func validSet(results []validator.Result) map[string]bool {
	s := make(map[string]bool, len(results))
	for _, r := range results {
		if r.OK {
			s[r.PatchName] = true
		}
	}
	return s
}
