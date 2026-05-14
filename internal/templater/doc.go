// Package templater provides Go text/template rendering for patch metadata.
//
// Templates receive a Data struct with the following fields:
//
//	{{.Patch.Name}}   - patch name from config
//	{{.Patch.Repo}}   - target repository path
//	{{.Patch.Tags}}   - slice of tags
//	{{.Status}}       - current state status (applied/pending/failed)
//	{{.Repo}}         - repo path recorded in state
//	{{.Date}}         - today's date (YYYY-MM-DD)
//	{{.Error}}        - last error message, if any
//
// Example usage:
//
//	results := templater.Render(patches, stateEntries,
//	    "[{{.Patch.Name}}] {{.Status}} on {{.Repo}} ({{.Date}})")
//	if templater.AnyFailed(results) { ... }
package templater
