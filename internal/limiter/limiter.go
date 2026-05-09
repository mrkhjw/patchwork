// Package limiter provides concurrency control for patch operations,
// allowing a maximum number of patches to run simultaneously across repos.
package limiter

import (
	"fmt"
	"sync"
)

// Result holds the outcome of a single limited operation.
type Result struct {
	Name  string
	Error error
}

// AnyFailed returns true if any result contains a non-nil error.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if r.Error != nil {
			return true
		}
	}
	return false
}

// Task is a named unit of work to be executed under the concurrency limit.
type Task struct {
	Name string
	Fn   func() error
}

// Run executes all tasks with at most maxConcurrent goroutines running at once.
// It returns one Result per task preserving submission order.
func Run(tasks []Task, maxConcurrent int) []Result {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}

	results := make([]Result, len(tasks))
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i, t := range tasks {
		wg.Add(1)
		go func(idx int, task Task) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			var err error
			if task.Fn == nil {
				err = fmt.Errorf("task %q has no function", task.Name)
			} else {
				err = task.Fn()
			}
			results[idx] = Result{Name: task.Name, Error: err}
		}(i, t)
	}

	wg.Wait()
	return results
}
