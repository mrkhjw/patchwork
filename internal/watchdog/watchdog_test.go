package watchdog_test

import (
	"testing"
	"time"

	"github.com/patchwork/internal/state"
	"github.com/patchwork/internal/watchdog"
)

var baseTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func entry(name, repo, status string, age time.Duration) state.Entry {
	return state.Entry{
		PatchName: name,
		Repo:      repo,
		Status:    status,
		UpdatedAt: baseTime.Add(-age),
	}
}

func TestWatch_NoAlerts_FreshEntries(t *testing.T) {
	entries := []state.Entry{
		entry("p1", "repo-a", "pending", 1*time.Hour),
	}
	alerts := watchdog.Watch(entries, baseTime, watchdog.DefaultPolicy())
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestWatch_WarnAlert_OldPending(t *testing.T) {
	entries := []state.Entry{
		entry("p1", "repo-a", "pending", 30*time.Hour),
	}
	alerts := watchdog.Watch(entries, baseTime, watchdog.DefaultPolicy())
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != watchdog.SeverityWarn {
		t.Errorf("expected warn, got %s", alerts[0].Severity)
	}
}

func TestWatch_CritAlert_VeryOldFailed(t *testing.T) {
	entries := []state.Entry{
		entry("p2", "repo-b", "failed", 96*time.Hour),
	}
	alerts := watchdog.Watch(entries, baseTime, watchdog.DefaultPolicy())
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != watchdog.SeverityCrit {
		t.Errorf("expected critical, got %s", alerts[0].Severity)
	}
}

func TestWatch_SkipsAppliedStatus(t *testing.T) {
	entries := []state.Entry{
		entry("p3", "repo-c", "applied", 200*time.Hour),
	}
	alerts := watchdog.Watch(entries, baseTime, watchdog.DefaultPolicy())
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for applied status, got %d", len(alerts))
	}
}

func TestAnyFailed_True(t *testing.T) {
	alerts := []watchdog.Alert{
		{Severity: watchdog.SeverityWarn},
		{Severity: watchdog.SeverityCrit},
	}
	if !watchdog.AnyFailed(alerts) {
		t.Error("expected AnyFailed to be true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	alerts := []watchdog.Alert{
		{Severity: watchdog.SeverityWarn},
	}
	if watchdog.AnyFailed(alerts) {
		t.Error("expected AnyFailed to be false")
	}
}

func TestWatch_AlertMessage_ContainsPatchName(t *testing.T) {
	entries := []state.Entry{
		entry("my-patch", "repo-x", "pending", 48*time.Hour),
	}
	alerts := watchdog.Watch(entries, baseTime, watchdog.DefaultPolicy())
	if len(alerts) == 0 {
		t.Fatal("expected at least one alert")
	}
	if alerts[0].PatchName != "my-patch" {
		t.Errorf("expected patch name my-patch, got %s", alerts[0].PatchName)
	}
}
