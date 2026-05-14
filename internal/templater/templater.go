// Package templater renders patch metadata into text templates,
// useful for generating commit messages, PR descriptions, or notifications.
package templater

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/state"
)

// Data holds the values available inside a template.
type Data struct {
	Patch  config.Patch
	Status string
	Repo   string
	Date   string
	Error  string
}

// Result is the outcome of rendering a single patch.
type Result struct {
	PatchName string
	Output    string
	Err       error
}

// Render executes tmplStr for each patch, enriching it with state info.
func Render(patches []config.Patch, entries map[string]state.Entry, tmplStr string) []Result {
	tmpl, err := template.New("patch").Parse(tmplStr)
	if err != nil {
		results := make([]Result, len(patches))
		for i, p := range patches {
			results[i] = Result{PatchName: p.Name, Err: fmt.Errorf("parse template: %w", err)}
		}
		return results
	}

	results := make([]Result, 0, len(patches))
	for _, p := range patches {
		d := Data{
			Patch: p,
			Date:  time.Now().Format(time.DateOnly),
		}
		if e, ok := entries[p.Name]; ok {
			d.Status = e.Status
			d.Repo = e.Repo
			d.Error = e.Error
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, d); err != nil {
			results = append(results, Result{PatchName: p.Name, Err: fmt.Errorf("execute template: %w", err)})
			continue
		}
		results = append(results, Result{PatchName: p.Name, Output: buf.String()})
	}
	return results
}

// AnyFailed returns true if any result carries an error.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
