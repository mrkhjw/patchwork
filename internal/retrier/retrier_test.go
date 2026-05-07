package retrier

import (
	"errors"
	"testing"
	"time"
)

var noSleep = func(time.Duration) {}

func TestRun_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	r := Run(DefaultPolicy(), func() error {
		calls++
		return nil
	}, noSleep)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	if r.Attempts != 1 {
		t.Fatalf("expected Attempts=1, got %d", r.Attempts)
	}
}

func TestRun_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	r := Run(DefaultPolicy(), func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	}, noSleep)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", r.Attempts)
	}
}

func TestRun_ExhaustsAttempts(t *testing.T) {
	p := Policy{MaxAttempts: 3, Delay: 0, Backoff: 1.0}
	calls := 0
	r := Run(p, func() error {
		calls++
		return errors.New("always fails")
	}, noSleep)
	if r.Err == nil {
		t.Fatal("expected error")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_BackoffCallsSleep(t *testing.T) {
	p := Policy{MaxAttempts: 3, Delay: 10 * time.Millisecond, Backoff: 2.0}
	var slept []time.Duration
	recorder := func(d time.Duration) { slept = append(slept, d) }

	Run(p, func() error { return errors.New("fail") }, recorder)

	if len(slept) != 2 {
		t.Fatalf("expected 2 sleeps, got %d", len(slept))
	}
	if slept[1] != 2*slept[0] {
		t.Fatalf("expected backoff doubling: %v -> %v", slept[0], slept[1])
	}
}

func TestRunSmart_StopsOnPermanent(t *testing.T) {
	p := Policy{MaxAttempts: 5, Delay: 0, Backoff: 1.0}
	calls := 0
	r := RunSmart(p, func() error {
		calls++
		return Permanent{Err: errors.New("fatal")}
	}, noSleep)
	if r.Err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call on permanent error, got %d", calls)
	}
}

func TestIsPermanent(t *testing.T) {
	perm := Permanent{Err: errors.New("oops")}
	if !IsPermanent(perm) {
		t.Fatal("expected IsPermanent=true")
	}
	if IsPermanent(errors.New("regular")) {
		t.Fatal("expected IsPermanent=false for regular error")
	}
}

func TestRun_ZeroMaxAttemptsDefaultsToOne(t *testing.T) {
	p := Policy{MaxAttempts: 0, Delay: 0, Backoff: 1.0}
	calls := 0
	Run(p, func() error { calls++; return errors.New("x") }, noSleep)
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
