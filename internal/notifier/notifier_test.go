package notifier_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/patchwork/internal/notifier"
	"github.com/yourorg/patchwork/internal/reporter"
)

func sampleEntries() []reporter.Entry {
	return []reporter.Entry{
		{Patch: "fix-a", Repo: "repo1", Status: "applied"},
		{Patch: "fix-b", Repo: "repo2", Status: "skipped"},
		{Patch: "fix-c", Repo: "repo3", Status: "failed"},
	}
}

func TestNotify_FileTarget(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "out.log")

	targets := []notifier.Target{
		{Kind: "file", Destination: logPath},
	}

	if err := notifier.Notify(targets, sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if !strings.Contains(string(data), "patchwork run complete") {
		t.Errorf("expected summary in file, got: %s", string(data))
	}
}

func TestNotify_UnknownKind(t *testing.T) {
	targets := []notifier.Target{
		{Kind: "webhook"},
	}
	err := notifier.Notify(targets, sampleEntries())
	if err == nil {
		t.Fatal("expected error for unknown kind, got nil")
	}
	if !strings.Contains(err.Error(), "unknown kind") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNotify_MultipleTargets_ContinuesOnError(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "out.log")

	targets := []notifier.Target{
		{Kind: "bad-kind"},
		{Kind: "file", Destination: logPath},
	}

	err := notifier.Notify(targets, sampleEntries())
	if err == nil {
		t.Fatal("expected error from bad target")
	}

	// second target (file) should still have been written
	if _, statErr := os.Stat(logPath); os.IsNotExist(statErr) {
		t.Error("file target was not written despite earlier error")
	}
}

func TestAnyFailed_True(t *testing.T) {
	if !notifier.AnyFailed(sampleEntries()) {
		t.Error("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	entries := []reporter.Entry{
		{Patch: "p", Repo: "r", Status: "applied"},
		{Patch: "q", Repo: "r", Status: "skipped"},
	}
	if notifier.AnyFailed(entries) {
		t.Error("expected AnyFailed to return false")
	}
}
