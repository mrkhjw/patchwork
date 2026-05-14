package profiler_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/patchwork/internal/profiler"
)

func TestFormat_ShowsErrorStatus(t *testing.T) {
	p := profiler.New()
	_ = p.Track("broken", func() error { return errors.New("disk full") })
	out := profiler.Format(p)
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected error message in format output:\n%s", out)
	}
}

func TestFormat_ShowsOkStatus(t *testing.T) {
	p := profiler.New()
	_ = p.Track("healthy", func() error { return nil })
	out := profiler.Format(p)
	if !strings.Contains(out, "ok") {
		t.Errorf("expected 'ok' status in format output:\n%s", out)
	}
}

func TestFormat_ShowsOperationCount(t *testing.T) {
	p := profiler.New()
	_ = p.Track("a", func() error { return nil })
	_ = p.Track("b", func() error { return nil })
	_ = p.Track("c", func() error { return nil })
	out := profiler.Format(p)
	if !strings.Contains(out, "3 operations") {
		t.Errorf("expected '3 operations' in output:\n%s", out)
	}
}

func TestFormat_MultipleRecords_AllNamed(t *testing.T) {
	p := profiler.New()
	names := []string{"alpha", "beta", "gamma"}
	for _, n := range names {
		name := n
		_ = p.Track(name, func() error { return nil })
	}
	out := profiler.Format(p)
	for _, n := range names {
		if !strings.Contains(out, n) {
			t.Errorf("expected name %q in output:\n%s", n, out)
		}
	}
}
