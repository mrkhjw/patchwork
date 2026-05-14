package profiler_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/patchwork/internal/profiler"
)

func TestTrack_RecordsDuration(t *testing.T) {
	p := profiler.New()
	_ = p.Track("op1", func() error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})
	if len(p.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(p.Records))
	}
	if p.Records[0].Duration < 5*time.Millisecond {
		t.Errorf("expected duration >= 5ms, got %s", p.Records[0].Duration)
	}
	if p.Records[0].Name != "op1" {
		t.Errorf("unexpected name: %s", p.Records[0].Name)
	}
}

func TestTrack_RecordsError(t *testing.T) {
	p := profiler.New()
	sentinel := errors.New("boom")
	err := p.Track("failing", func() error { return sentinel })
	if err != sentinel {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if p.Records[0].Err != sentinel {
		t.Errorf("error not stored in record")
	}
}

func TestAnyFailed_False(t *testing.T) {
	p := profiler.New()
	_ = p.Track("ok", func() error { return nil })
	if p.AnyFailed() {
		t.Error("expected AnyFailed to be false")
	}
}

func TestAnyFailed_True(t *testing.T) {
	p := profiler.New()
	_ = p.Track("ok", func() error { return nil })
	_ = p.Track("bad", func() error { return errors.New("fail") })
	if !p.AnyFailed() {
		t.Error("expected AnyFailed to be true")
	}
}

func TestSlowest_ReturnsSorted(t *testing.T) {
	p := profiler.New()
	_ = p.Track("fast", func() error { time.Sleep(2 * time.Millisecond); return nil })
	_ = p.Track("slow", func() error { time.Sleep(10 * time.Millisecond); return nil })
	_ = p.Track("mid", func() error { time.Sleep(5 * time.Millisecond); return nil })

	top := p.Slowest(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 records, got %d", len(top))
	}
	if top[0].Duration < top[1].Duration {
		t.Errorf("records not sorted descending: %s < %s", top[0].Duration, top[1].Duration)
	}
}

func TestSlowest_ClampsToLen(t *testing.T) {
	p := profiler.New()
	_ = p.Track("a", func() error { return nil })
	result := p.Slowest(100)
	if len(result) != 1 {
		t.Errorf("expected 1, got %d", len(result))
	}
}

func TestTotalDuration(t *testing.T) {
	p := profiler.New()
	_ = p.Track("a", func() error { time.Sleep(5 * time.Millisecond); return nil })
	_ = p.Track("b", func() error { time.Sleep(5 * time.Millisecond); return nil })
	if p.TotalDuration() < 10*time.Millisecond {
		t.Errorf("total duration too small: %s", p.TotalDuration())
	}
}

func TestFormat_ContainsName(t *testing.T) {
	p := profiler.New()
	_ = p.Track("my-operation", func() error { return nil })
	out := profiler.Format(p)
	if !strings.Contains(out, "my-operation") {
		t.Errorf("format output missing operation name:\n%s", out)
	}
}

func TestFormat_EmptyProfile(t *testing.T) {
	p := profiler.New()
	out := profiler.Format(p)
	if !strings.Contains(out, "no records") {
		t.Errorf("expected 'no records' message, got:\n%s", out)
	}
}
