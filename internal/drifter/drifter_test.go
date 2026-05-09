package drifter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork/internal/drifter"
	"github.com/patchwork/internal/snapshot"
	"github.com/patchwork/internal/state"
)

func makeSnap(entries []snapshot.Entry) snapshot.Snapshot {
	return snapshot.Snapshot{TakenAt: time.Now().UTC(), Entries: entries}
}

func makeState(entries []state.Entry) *state.State {
	st := state.New()
	for _, e := range entries {
		st.Upsert(e)
	}
	return st
}

func TestDetect_NoDrift(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{PatchName: "fix-a", Repo: "/repos/a", Status: "applied"},
	})
	st := makeState([]state.Entry{
		{PatchName: "fix-a", Repo: "/repos/a", Status: "applied"},
	})
	r := drifter.Detect(snap, st)
	if r.AnyDrifted() {
		t.Fatalf("expected no drift, got %d alerts", len(r.Alerts))
	}
}

func TestDetect_StatusChanged(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{PatchName: "fix-b", Repo: "/repos/b", Status: "applied"},
	})
	st := makeState([]state.Entry{
		{PatchName: "fix-b", Repo: "/repos/b", Status: "pending"},
	})
	r := drifter.Detect(snap, st)
	if !r.AnyDrifted() {
		t.Fatal("expected drift")
	}
	if r.Alerts[0].Field != "status" {
		t.Errorf("expected field=status, got %q", r.Alerts[0].Field)
	}
}

func TestDetect_MissingEntry(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{PatchName: "fix-c", Repo: "/repos/c", Status: "applied"},
	})
	st := makeState(nil)
	r := drifter.Detect(snap, st)
	if len(r.Alerts) != 1 || r.Alerts[0].Field != "existence" {
		t.Errorf("expected existence alert, got %+v", r.Alerts)
	}
}

func TestDetect_RepoChanged(t *testing.T) {
	snap := makeSnap([]snapshot.Entry{
		{PatchName: "fix-d", Repo: "/repos/old", Status: "applied"},
	})
	st := makeState([]state.Entry{
		{PatchName: "fix-d", Repo: "/repos/new", Status: "applied"},
	})
	r := drifter.Detect(snap, st)
	if !r.AnyDrifted() {
		t.Fatal("expected drift on repo field")
	}
	if r.Alerts[0].Field != "repo" {
		t.Errorf("expected field=repo, got %q", r.Alerts[0].Field)
	}
}

func TestFormat_NoDrift(t *testing.T) {
	r := drifter.Result{}
	out := drifter.Format(r)
	if out != "no drift detected" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithAlerts(t *testing.T) {
	r := drifter.Result{
		Alerts: []drifter.Alert{
			{PatchName: "fix-e", Repo: "/repos/e", Field: "status", Expected: "applied", Actual: "pending"},
		},
	}
	out := drifter.Format(r)
	if !strings.Contains(out, "fix-e") || !strings.Contains(out, "status") {
		t.Errorf("format missing expected content: %q", out)
	}
}
