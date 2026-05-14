package templater_test

import (
	"strings"
	"testing"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/state"
	"github.com/yourorg/patchwork/internal/templater"
)

func samplePatches() []config.Patch {
	return []config.Patch{
		{Name: "fix-auth", Repo: "/repos/api", Tags: []string{"security"}},
		{Name: "add-logging", Repo: "/repos/svc", Tags: []string{"ops"}},
	}
}

func sampleEntries() map[string]state.Entry {
	return map[string]state.Entry{
		"fix-auth":    {Status: "applied", Repo: "/repos/api"},
		"add-logging": {Status: "pending", Repo: "/repos/svc"},
	}
}

func TestRender_BasicTemplate(t *testing.T) {
	results := templater.Render(samplePatches(), sampleEntries(), "{{.Patch.Name}}:{{.Status}}")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Output != "fix-auth:applied" {
		t.Errorf("unexpected output: %q", results[0].Output)
	}
	if results[1].Output != "add-logging:pending" {
		t.Errorf("unexpected output: %q", results[1].Output)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	results := templater.Render(samplePatches(), sampleEntries(), "{{.Unclosed")
	for _, r := range results {
		if r.Err == nil {
			t.Errorf("expected error for patch %q", r.PatchName)
		}
	}
}

func TestRender_MissingStateEntry(t *testing.T) {
	entries := map[string]state.Entry{} // empty
	results := templater.Render(samplePatches(), entries, "{{.Patch.Name}}:{{.Status}}")
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error: %v", r.Err)
		}
		if !strings.HasSuffix(r.Output, ":") {
			t.Errorf("expected empty status, got %q", r.Output)
		}
	}
}

func TestRender_ContainsDate(t *testing.T) {
	results := templater.Render(samplePatches()[:1], sampleEntries(), "date={{.Date}}")
	if results[0].Err != nil {
		t.Fatal(results[0].Err)
	}
	if !strings.HasPrefix(results[0].Output, "date=20") {
		t.Errorf("expected date prefix, got %q", results[0].Output)
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := templater.Render(samplePatches(), sampleEntries(), "{{.Patch.Name}}")
	if templater.AnyFailed(results) {
		t.Error("expected no failures")
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := templater.Render(samplePatches(), sampleEntries(), "{{.Bad")
	if !templater.AnyFailed(results) {
		t.Error("expected failures")
	}
}
