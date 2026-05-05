package patch

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initRepo creates a temporary git repository with a single committed file.
func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@example.com")
	run("git", "config", "user.name", "Test")

	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "init")

	return dir
}

// writePatch writes a unified diff patch file and returns its path.
func writePatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.patch")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

const simplePatch = `--- a/hello.txt
+++ b/hello.txt
@@ -1 +1,2 @@
 hello
+world
`

func TestApply_Success(t *testing.T) {
	repo := initRepo(t)
	pf := writePatch(t, simplePatch)

	p := Patch{Name: "add-world", File: pf, RepoDir: repo}
	if err := Apply(p); err != nil {
		t.Fatalf("Apply returned unexpected error: %v", err)
	}
}

func TestApply_MissingFile(t *testing.T) {
	repo := initRepo(t)
	p := Patch{Name: "missing", File: "/nonexistent/file.patch", RepoDir: repo}
	if err := Apply(p); err == nil {
		t.Fatal("expected error for missing patch file, got nil")
	}
}

func TestCheck_Pending(t *testing.T) {
	repo := initRepo(t)
	pf := writePatch(t, simplePatch)

	p := Patch{Name: "add-world", File: pf, RepoDir: repo}
	if got := Check(p); got != StatusPending {
		t.Fatalf("expected %s, got %s", StatusPending, got)
	}
}

func TestCheck_Applied(t *testing.T) {
	repo := initRepo(t)
	pf := writePatch(t, simplePatch)

	p := Patch{Name: "add-world", File: pf, RepoDir: repo}
	if err := Apply(p); err != nil {
		t.Fatalf("setup Apply failed: %v", err)
	}
	if got := Check(p); got != StatusApplied {
		t.Fatalf("expected %s, got %s", StatusApplied, got)
	}
}
