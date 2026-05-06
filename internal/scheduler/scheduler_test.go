package scheduler_test

import (
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/scheduler"
)

func named(names ...string) []config.Patch {
	patches := make([]config.Patch, len(names))
	for i, n := range names {
		patches[i] = config.Patch{Name: n}
	}
	return patches
}

func TestOrder_NoDeps(t *testing.T) {
	patches := named("a", "b", "c")
	got, err := scheduler.Order(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 patches, got %d", len(got))
	}
}

func TestOrder_LinearDeps(t *testing.T) {
	patches := []config.Patch{
		{Name: "c", DependsOn: []string{"b"}},
		{Name: "b", DependsOn: []string{"a"}},
		{Name: "a"},
	}
	got, err := scheduler.Order(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	order := make(map[string]int, len(got))
	for i, p := range got {
		order[p.Name] = i
	}
	if order["a"] >= order["b"] || order["b"] >= order["c"] {
		t.Errorf("unexpected order: a=%d b=%d c=%d", order["a"], order["b"], order["c"])
	}
}

func TestOrder_DetectsCycle(t *testing.T) {
	patches := []config.Patch{
		{Name: "a", DependsOn: []string{"b"}},
		{Name: "b", DependsOn: []string{"a"}},
	}
	_, err := scheduler.Order(patches)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
	var cycleErr scheduler.ErrCycle
	if _, ok := err.(scheduler.ErrCycle); !ok {
		t.Errorf("expected ErrCycle, got %T: %v", err, err)
	}
	_ = cycleErr
}

func TestOrder_UnknownDep(t *testing.T) {
	patches := []config.Patch{
		{Name: "a", DependsOn: []string{"missing"}},
	}
	_, err := scheduler.Order(patches)
	if err == nil {
		t.Fatal("expected error for unknown dependency, got nil")
	}
}

func TestOrder_DiamondDeps(t *testing.T) {
	patches := []config.Patch{
		{Name: "d", DependsOn: []string{"b", "c"}},
		{Name: "b", DependsOn: []string{"a"}},
		{Name: "c", DependsOn: []string{"a"}},
		{Name: "a"},
	}
	got, err := scheduler.Order(patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 patches, got %d", len(got))
	}
	pos := make(map[string]int)
	for i, p := range got {
		pos[p.Name] = i
	}
	if pos["a"] >= pos["b"] || pos["a"] >= pos["c"] || pos["b"] >= pos["d"] || pos["c"] >= pos["d"] {
		t.Errorf("diamond order violated: %v", pos)
	}
}
