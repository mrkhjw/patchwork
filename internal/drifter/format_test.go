package drifter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork/internal/drifter"
)

func sampleAlerts() []drifter.Alert {
	return []drifter.Alert{
		{
			PatchName:  "patch-x",
			Repo:       "/repos/x",
			Field:      "status",
			Expected:   "applied",
			Actual:     "failed",
			DetectedAt: time.Now().UTC(),
		},
		{
			PatchName:  "patch-y",
			Repo:       "/repos/y",
			Field:      "existence",
			Expected:   "present",
			Actual:     "missing",
			DetectedAt: time.Now().UTC(),
		},
	}
}

func TestFormat_AlertCount(t *testing.T) {
	r := drifter.Result{Alerts: sampleAlerts()}
	out := drifter.Format(r)
	if !strings.Contains(out, "2 drift alert") {
		t.Errorf("expected alert count in output, got: %q", out)
	}
}

func TestFormat_ContainsRepo(t *testing.T) {
	r := drifter.Result{Alerts: sampleAlerts()}
	out := drifter.Format(r)
	if !strings.Contains(out, "/repos/x") {
		t.Errorf("expected repo path in output, got: %q", out)
	}
}

func TestFormat_ContainsExpectedAndActual(t *testing.T) {
	r := drifter.Result{Alerts: sampleAlerts()}
	out := drifter.Format(r)
	if !strings.Contains(out, "applied") || !strings.Contains(out, "failed") {
		t.Errorf("expected expected/actual values in output, got: %q", out)
	}
}

func TestFormat_MultipleEntries(t *testing.T) {
	r := drifter.Result{Alerts: sampleAlerts()}
	out := drifter.Format(r)
	if !strings.Contains(out, "patch-x") || !strings.Contains(out, "patch-y") {
		t.Errorf("expected both patches in output, got: %q", out)
	}
}
