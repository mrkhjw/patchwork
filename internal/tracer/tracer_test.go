package tracer

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestStart_RecordsSpan(t *testing.T) {
	tr := New()
	done := tr.Start("fix-auth", "repo-a")
	time.Sleep(2 * time.Millisecond)
	done(nil)

	spans := tr.Spans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	s := spans[0]
	if s.PatchName != "fix-auth" || s.Repo != "repo-a" {
		t.Errorf("unexpected span fields: %+v", s)
	}
	if s.Duration < 2*time.Millisecond {
		t.Errorf("duration too short: %v", s.Duration)
	}
	if s.Err != nil {
		t.Errorf("expected nil error")
	}
}

func TestStart_RecordsError(t *testing.T) {
	tr := New()
	err := errors.New("patch failed")
	done := tr.Start("fix-db", "repo-b")
	done(err)

	if !tr.AnyFailed() {
		t.Error("expected AnyFailed to be true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	tr := New()
	tr.Start("p1", "r1")(nil)
	tr.Start("p2", "r2")(nil)
	if tr.AnyFailed() {
		t.Error("expected AnyFailed to be false")
	}
}

func TestSlowest_ReturnsSorted(t *testing.T) {
	tr := New()
	// inject spans with known durations
	tr.spans = []Span{
		{PatchName: "fast", Duration: 1 * time.Millisecond},
		{PatchName: "slow", Duration: 50 * time.Millisecond},
		{PatchName: "medium", Duration: 10 * time.Millisecond},
	}

	top := tr.Slowest(2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
	if top[0].PatchName != "slow" {
		t.Errorf("expected slowest first, got %s", top[0].PatchName)
	}
	if top[1].PatchName != "medium" {
		t.Errorf("expected medium second, got %s", top[1].PatchName)
	}
}

func TestSlowest_ClampsToLen(t *testing.T) {
	tr := New()
	tr.spans = []Span{{PatchName: "only", Duration: 5 * time.Millisecond}}
	if got := tr.Slowest(10); len(got) != 1 {
		t.Errorf("expected 1, got %d", len(got))
	}
}

func TestWrite_ContainsFields(t *testing.T) {
	tr := New()
	tr.spans = []Span{
		{PatchName: "my-patch", Repo: "my-repo", Duration: 3 * time.Millisecond, Err: nil},
		{PatchName: "bad-patch", Repo: "other-repo", Duration: 1 * time.Millisecond, Err: errors.New("boom")},
	}

	var buf bytes.Buffer
	tr.Write(&buf)
	out := buf.String()

	for _, want := range []string{"my-patch", "my-repo", "bad-patch", "boom", "PATCH", "DURATION"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}
