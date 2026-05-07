package digester_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork/internal/config"
	"github.com/patchwork/internal/digester"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "patch-*.diff")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDigest_Success(t *testing.T) {
	path := writeTempFile(t, "--- a/foo\n+++ b/foo\n@@ -1 +1 @@\n-old\n+new\n")
	patches := []config.Patch{{Name: "fix-foo", Path: path}}

	results := digester.Digest(patches)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if !r.OK() {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if len(r.Digest) != 64 {
		t.Errorf("expected 64-char hex digest, got %q", r.Digest)
	}
}

func TestDigest_MissingFile(t *testing.T) {
	patches := []config.Patch{{Name: "ghost", Path: filepath.Join(t.TempDir(), "missing.diff")}}
	results := digester.Digest(patches)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].OK() {
		t.Error("expected error for missing file")
	}
}

func TestAnyFailed_True(t *testing.T) {
	results := []digester.Result{
		{PatchName: "ok", Digest: "abc"},
		{PatchName: "bad", Err: os.ErrNotExist},
	}
	if !digester.AnyFailed(results) {
		t.Error("expected AnyFailed to return true")
	}
}

func TestAnyFailed_False(t *testing.T) {
	results := []digester.Result{
		{PatchName: "ok", Digest: "abc"},
	}
	if digester.AnyFailed(results) {
		t.Error("expected AnyFailed to return false")
	}
}

func TestVerify_Match(t *testing.T) {
	path := writeTempFile(t, "hello world")
	results := digester.Digest([]config.Patch{{Name: "p", Path: path}})
	if err := digester.Verify(path, results[0].Digest); err != nil {
		t.Errorf("expected verify to pass: %v", err)
	}
}

func TestVerify_Mismatch(t *testing.T) {
	path := writeTempFile(t, "hello world")
	if err := digester.Verify(path, "0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Error("expected verify to fail on mismatch")
	}
}

func TestVerify_MissingFile(t *testing.T) {
	if err := digester.Verify(filepath.Join(t.TempDir(), "nope.diff"), "abc"); err == nil {
		t.Error("expected error for missing file")
	}
}
