package limiter_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/patchwork/internal/limiter"
)

func TestRun_AllSucceed(t *testing.T) {
	tasks := []limiter.Task{
		{Name: "a", Fn: func() error { return nil }},
		{Name: "b", Fn: func() error { return nil }},
		{Name: "c", Fn: func() error { return nil }},
	}
	results := limiter.Run(tasks, 2)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Error != nil {
			t.Errorf("unexpected error for %q: %v", r.Name, r.Error)
		}
	}
}

func TestRun_PartialFailure(t *testing.T) {
	errBoom := errors.New("boom")
	tasks := []limiter.Task{
		{Name: "ok", Fn: func() error { return nil }},
		{Name: "fail", Fn: func() error { return errBoom }},
	}
	results := limiter.Run(tasks, 2)
	if results[0].Error != nil {
		t.Errorf("expected no error for 'ok', got %v", results[0].Error)
	}
	if !errors.Is(results[1].Error, errBoom) {
		t.Errorf("expected boom error for 'fail', got %v", results[1].Error)
	}
}

func TestRun_NilFn_ReturnsError(t *testing.T) {
	tasks := []limiter.Task{
		{Name: "nilTask", Fn: nil},
	}
	results := limiter.Run(tasks, 1)
	if results[0].Error == nil {
		t.Error("expected error for nil Fn, got nil")
	}
}

func TestRun_RespectsConcurrencyLimit(t *testing.T) {
	var active int64
	var peak int64

	tasks := make([]limiter.Task, 10)
	for i := range tasks {
		tasks[i] = limiter.Task{
			Name: "t",
			Fn: func() error {
				cur := atomic.AddInt64(&active, 1)
				for {
					old := atomic.LoadInt64(&peak)
					if cur <= old || atomic.CompareAndSwapInt64(&peak, old, cur) {
						break
					}
				}
				time.Sleep(5 * time.Millisecond)
				atomic.AddInt64(&active, -1)
				return nil
			},
		}
	}

	limiter.Run(tasks, 3)
	if peak > 3 {
		t.Errorf("concurrency exceeded limit: peak=%d, want<=3", peak)
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := []limiter.Result{
		{Name: "a", Error: nil},
		{Name: "b", Error: errors.New("fail")},
	}
	if !limiter.AnyFailed(results) {
		t.Error("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := []limiter.Result{
		{Name: "a", Error: nil},
		{Name: "b", Error: nil},
	}
	if limiter.AnyFailed(results) {
		t.Error("expected AnyFailed to return false")
	}
}
