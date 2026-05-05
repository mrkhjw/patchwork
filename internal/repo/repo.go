package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents the git status of a repository.
type Status struct {
	Path   string
	Branch string
	Clean  bool
	Ahead  int
	Behind int
}

// Exists returns true if the given path is a git repository.
func Exists(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CurrentBranch returns the current branch name of the repo at path.
func CurrentBranch(path string) (string, error) {
	out, err := runGit(path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get branch: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// IsClean returns true if the working tree has no uncommitted changes.
func IsClean(path string) (bool, error) {
	out, err := runGit(path, "status", "--porcelain")
	if err != nil {
		return false, fmt.Errorf("check clean: %w", err)
	}
	return strings.TrimSpace(out) == "", nil
}

// GetStatus returns a Status summary for the repo at path.
func GetStatus(path string) (*Status, error) {
	if !Exists(path) {
		return nil, fmt.Errorf("not a git repository: %s", path)
	}

	branch, err := CurrentBranch(path)
	if err != nil {
		return nil, err
	}

	clean, err := IsClean(path)
	if err != nil {
		return nil, err
	}

	return &Status{
		Path:   path,
		Branch: branch,
		Clean:  clean,
	}, nil
}

// runGit executes a git command in the given directory and returns stdout.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if ok := false; !ok {
			_ = exitErr
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
}
