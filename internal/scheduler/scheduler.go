// Package scheduler provides ordered execution planning for patches,
// respecting declared dependencies between patches.
package scheduler

import (
	"fmt"

	"github.com/patchwork/internal/config"
)

// ErrCycle is returned when a dependency cycle is detected.
type ErrCycle struct {
	Patch string
}

func (e ErrCycle) Error() string {
	return fmt.Sprintf("cycle detected involving patch %q", e.Patch)
}

// Order returns patches in dependency-resolved order using a topological sort.
// Patches with no declared dependencies are returned in their original order.
func Order(patches []config.Patch) ([]config.Patch, error) {
	index := make(map[string]config.Patch, len(patches))
	for _, p := range patches {
		index[p.Name] = p
	}

	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)

	state := make(map[string]int, len(patches))
	var sorted []config.Patch

	var visit func(name string) error
	visit = func(name string) error {
		switch state[name] {
		case visited:
			return nil
		case visiting:
			return ErrCycle{Patch: name}
		}
		state[name] = visiting
		p, ok := index[name]
		if !ok {
			return fmt.Errorf("unknown patch dependency %q", name)
		}
		for _, dep := range p.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[name] = visited
		sorted = append(sorted, p)
		return nil
	}

	for _, p := range patches {
		if err := visit(p.Name); err != nil {
			return nil, err
		}
	}
	return sorted, nil
}
