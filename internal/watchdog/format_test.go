package watchdog_test

import (
	"strings"
	"testing"
	"time"

	"github.com/patchwork/internal/watchdog"
)

func TestFormat_EmptyAlerts(t *testing.T) {
	out := watchdog.Format(nil)
	if !strings.Contains(out, "no alerts") {
		t.Errorf("expected 'no alerts' in output, got: %s", out)
	}
}

func TestFormat_ContainsHeaders(t *testing.T) {
	alerts := []watchdog.Alert{
		{
			PatchName: "fix-auth",
			Repo:      "api-service",
			Status:    "pending",
			Age:       30 * time.Hour,
			Severity:  watchdog.SeverityWarn,
		},
	}
	out := watchdog.Format(alerts)
	for _, hdr := range []string{"SEVERITY", "PATCH", "REPO", "STATUS", "AGE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormat_ContainsAlertData(t *testing.T) {
	alerts := []watchdog.Alert{
		{
			PatchName: "db-migrate",
			Repo:      "core",
			Status:    "failed",
			Age:       80 * time.Hour,
			Severity:  watchdog.SeverityCrit,
		},
	}
	out := watchdog.Format(alerts)
	if !strings.Contains(out, "db-migrate") {
		t.Errorf("expected patch name in output")
	}
	if !strings.Contains(out, "critical") {
		t.Errorf("expected severity in output")
	}
	if !strings.Contains(out, "failed") {
		t.Errorf("expected status in output")
	}
}

func TestFormat_AlertCount(t *testing.T) {
	alerts := []watchdog.Alert{
		{PatchName: "a", Severity: watchdog.SeverityWarn, Age: 25 * time.Hour},
		{PatchName: "b", Severity: watchdog.SeverityCrit, Age: 80 * time.Hour},
	}
	out := watchdog.Format(alerts)
	if !strings.Contains(out, "2 alert(s)") {
		t.Errorf("expected alert count in output, got: %s", out)
	}
}
