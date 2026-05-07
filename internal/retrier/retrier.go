// Package retrier provides retry logic for patch operations that may
// transiently fail (e.g. due to locked files or network-backed repos).
package retrier

import (
	"errors"
	"fmt"
	"time"
)

// Policy controls how retries are attempted.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64 // multiplier applied to Delay after each attempt
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Backoff:     2.0,
	}
}

// Result holds the outcome of a retried operation.
type Result struct {
	Attempts int
	Err      error
}

// Run executes fn according to p, retrying on non-nil error.
// sleep is injectable for testing; pass nil to use time.Sleep.
func Run(p Policy, fn func() error, sleep func(time.Duration)) Result {
	if sleep == nil {
		sleep = time.Sleep
	}
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Backoff <= 0 {
		p.Backoff = 1.0
	}

	delay := p.Delay
	var lastErr error
	for i := 1; i <= p.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return Result{Attempts: i}
		} else {
			lastErr = err
		}
		if i < p.MaxAttempts {
			sleep(delay)
			delay = time.Duration(float64(delay) * p.Backoff)
		}
	}
	return Result{
		Attempts: p.MaxAttempts,
		Err:      fmt.Errorf("all %d attempts failed: %w", p.MaxAttempts, lastErr),
	}
}

// Permanent wraps an error to signal that retrying is futile.
type Permanent struct{ Err error }

func (p Permanent) Error() string { return p.Err.Error() }
func (p Permanent) Unwrap() error { return p.Err }

// IsPermanent reports whether err is a Permanent error.
func IsPermanent(err error) bool {
	var p Permanent
	return errors.As(err, &p)
}

// RunSmart is like Run but stops immediately on a Permanent error.
func RunSmart(p Policy, fn func() error, sleep func(time.Duration)) Result {
	if sleep == nil {
		sleep = time.Sleep
	}
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Backoff <= 0 {
		p.Backoff = 1.0
	}

	delay := p.Delay
	var lastErr error
	for i := 1; i <= p.MaxAttempts; i++ {
		err := fn()
		if err == nil {
			return Result{Attempts: i}
		}
		if IsPermanent(err) {
			return Result{Attempts: i, Err: err}
		}
		lastErr = err
		if i < p.MaxAttempts {
			sleep(delay)
			delay = time.Duration(float64(delay) * p.Backoff)
		}
	}
	return Result{
		Attempts: p.MaxAttempts,
		Err:      fmt.Errorf("all %d attempts failed: %w", p.MaxAttempts, lastErr),
	}
}
