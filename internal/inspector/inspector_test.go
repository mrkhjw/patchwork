package inspector_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/config"
	"github.com/yourorg/patchwork/internal/inspector"
	"github.com/yourorg/patchwork/internal/state"
)

func sampleConfig() *config.Config {
	return &config.Config{
		Patches: []config.Patch{
			{Name: "fix-auth", Repo: "/repos/auth", File: "fix-auth.patch", Tags: []string{"security"}},
			{Name: "bump-deps", Repo: "/repos/api", File: "bump-deps.patch", Tags: []string{}},
			{Name: "no-state", Repo: "/repos/core", File: "no-state.patch", Tags: []string{"infra"}},
		},
	}
}

func sampleState(t *testing.T) *state.State {
	t.Helper()
	st := state.New()
	st.Upsert(state.Entry{Name: "fix-auth", Repo: "/repos/auth", Status: "applied", AppliedAt: time.Now()})
	st.Upsert(state.Entry{Name: "bump-deps", Repo: "/repos/api", Status: "failed", AppliedAt: time.Time{}})
	return st
}

func TestInspect_StatusFromState(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Status != "applied" {
		t.Errorf("expected applied, got %s", results[0].Status)
	}
	if results[1].Status != "failed" {
		t.Errorf("expected failed, got %s", results[1].Status)
	}
}

func TestInspect_MissingStateEntry_WarnsPending(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	noState := results[2]
	if noState.Status != "pending" {
		t.Errorf("expected pending, got %s", noState.Status)
	}
	if len(noState.Warnings) == 0 {
		t.Error("expected warning for missing state entry")
	}
}

func TestInspect_NoTags_ProducesWarning(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	bumpDeps := results[1]
	found := false
	for _, w := range bumpDeps.Warnings {
		if strings.Contains(w, "no tags") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'no tags' warning for bump-deps")
	}
}

func TestInspect_FailedStatus_ProducesWarning(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	bumpDeps := results[1]
	found := false
	for _, w := range bumpDeps.Warnings {
		if strings.Contains(w, "failed state") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'failed state' warning")
	}
}

func TestAnyWarnings_True(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	if !inspector.AnyWarnings(results) {
		t.Error("expected AnyWarnings to return true")
	}
}

func TestAnyWarnings_False(t *testing.T) {
	clean := []inspector.Result{{Name: "ok", Status: "applied"}}
	if inspector.AnyWarnings(clean) {
		t.Error("expected AnyWarnings to return false")
	}
}

func TestFormat_ContainsHeaders(t *testing.T) {
	results := inspector.Inspect(sampleConfig(), sampleState(t))
	out := inspector.Format(results)
	for _, hdr := range []string{"NAME", "REPO", "STATUS", "WARNINGS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormat_EmptyResults(t *testing.T) {
	out := inspector.Format(nil)
	if !strings.Contains(out, "no patches") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
