package classifier_test

import (
	"strings"
	"testing"

	"github.com/patchwork/internal/classifier"
)

func sampleResults() []classifier.Result {
	return []classifier.Result{
		{PatchName: "fix-xss", Repo: "web", Tier: classifier.TierCritical, Reasons: []string{"critical tag: security"}},
		{PatchName: "add-index", Repo: "db", Tier: classifier.TierHigh, Reasons: []string{"high tag: schema"}},
		{PatchName: "update-copy", Repo: "web", Tier: classifier.TierLow, Reasons: nil},
	}
}

func TestFormat_ContainsHeader(t *testing.T) {
	out := classifier.Format(sampleResults())
	if !strings.Contains(out, "TIER") {
		t.Error("expected output to contain TIER header")
	}
}

func TestFormat_ContainsPatchName(t *testing.T) {
	out := classifier.Format(sampleResults())
	if !strings.Contains(out, "fix-xss") {
		t.Error("expected output to contain patch name")
	}
}

func TestFormat_ContainsTier(t *testing.T) {
	out := classifier.Format(sampleResults())
	if !strings.Contains(out, "critical") {
		t.Error("expected output to contain tier value")
	}
}

func TestFormat_EmptyResults(t *testing.T) {
	out := classifier.Format(nil)
	if !strings.Contains(out, "no patches") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormat_DashForNoReasons(t *testing.T) {
	out := classifier.Format(sampleResults())
	if !strings.Contains(out, "-") {
		t.Error("expected dash placeholder for empty reasons")
	}
}
