package drifter_test

import (
	"testing"
	"time"

	"github.com/patchwork/internal/drifter"
	"github.com/patchwork/internal/snapshot"
	"github.com/patchwork/internal/state"
)

// TestRoundtrip_SnapshotThenMutate captures a snapshot, mutates the state,
// and asserts that drift is correctly detected end-to-end.
func TestRoundtrip_SnapshotThenMutate(t *testing.T) {
	entries := []snapshot.Entry{
		{PatchName: "alpha", Repo: "/r/a", Status: "applied"},
		{PatchName: "beta", Repo: "/r/b", Status: "pending"},
	}
	snap := snapshot.Snapshot{TakenAt: time.Now().UTC(), Entries: entries}

	// Build state that matches snapshot exactly — no drift expected.
	st := state.New()
	st.Upsert(state.Entry{PatchName: "alpha", Repo: "/r/a", Status: "applied"})
	st.Upsert(state.Entry{PatchName: "beta", Repo: "/r/b", Status: "pending"})

	r := drifter.Detect(snap, st)
	if r.AnyDrifted() {
		t.Fatalf("expected no drift before mutation, got %d alerts", len(r.Alerts))
	}

	// Mutate: change beta status and remove alpha.
	st2 := state.New()
	st2.Upsert(state.Entry{PatchName: "beta", Repo: "/r/b", Status: "failed"})

	r2 := drifter.Detect(snap, st2)
	if !r2.AnyDrifted() {
		t.Fatal("expected drift after mutation")
	}

	fields := map[string]bool{}
	for _, a := range r2.Alerts {
		fields[a.Field] = true
	}
	if !fields["existence"] {
		t.Error("expected existence alert for removed alpha")
	}
	if !fields["status"] {
		t.Error("expected status alert for mutated beta")
	}
}
