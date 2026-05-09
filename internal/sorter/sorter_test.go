package sorter_test

import (
	"testing"

	"patchwork/internal/config"
	"patchwork/internal/sorter"
)

func sampleEntries() []sorter.Entry {
	return []sorter.Entry{
		{Patch: config.Patch{Name: "zebra-fix", Repo: "/repos/beta"}, Status: "pending"},
		{Patch: config.Patch{Name: "alpha-patch", Repo: "/repos/alpha"}, Status: "applied"},
		{Patch: config.Patch{Name: "middle-tweak", Repo: "/repos/beta"}, Status: "pending"},
		{Patch: config.Patch{Name: "core-update", Repo: "/repos/alpha"}, Status: "failed"},
	}
}

func TestSort_ByName(t *testing.T) {
	res := sorter.Sort(sampleEntries(), sorter.ByName)
	names := res.Names()
	expected := []string{"alpha-patch", "core-update", "middle-tweak", "zebra-fix"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("index %d: got %q, want %q", i, names[i], want)
		}
	}
}

func TestSort_ByRepo(t *testing.T) {
	res := sorter.Sort(sampleEntries(), sorter.ByRepo)
	// alpha entries first, then beta; within same repo sorted by name
	names := res.Names()
	if names[0] != "alpha-patch" || names[1] != "core-update" {
		t.Errorf("expected alpha repo entries first, got %v", names[:2])
	}
	if names[2] != "middle-tweak" || names[3] != "zebra-fix" {
		t.Errorf("expected beta repo entries last, got %v", names[2:])
	}
}

func TestSort_ByStatus(t *testing.T) {
	res := sorter.Sort(sampleEntries(), sorter.ByStatus)
	names := res.Names()
	// applied < failed < pending
	if names[0] != "alpha-patch" {
		t.Errorf("first should be applied entry, got %q", names[0])
	}
	if names[1] != "core-update" {
		t.Errorf("second should be failed entry, got %q", names[1])
	}
}

func TestSort_UnknownDimension_FallsBackToName(t *testing.T) {
	res := sorter.Sort(sampleEntries(), sorter.Dimension("unknown"))
	names := res.Names()
	if names[0] != "alpha-patch" {
		t.Errorf("expected alpha-patch first, got %q", names[0])
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	input := sampleEntries()
	first := input[0].Patch.Name
	sorter.Sort(input, sorter.ByName)
	if input[0].Patch.Name != first {
		t.Error("Sort mutated the original slice")
	}
}

func TestSort_EmptyEntries(t *testing.T) {
	res := sorter.Sort([]sorter.Entry{}, sorter.ByName)
	if len(res.Entries) != 0 {
		t.Errorf("expected empty result, got %d entries", len(res.Entries))
	}
}

func TestResult_Names(t *testing.T) {
	res := sorter.Sort(sampleEntries(), sorter.ByName)
	names := res.Names()
	if len(names) != 4 {
		t.Errorf("expected 4 names, got %d", len(names))
	}
}
