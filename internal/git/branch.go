package git

import (
	"os/exec"
	"strings"
)

// CurrentBranch returns the name of the current git branch.
func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// DefaultBranch tries to detect the default branch (main or master).
func DefaultBranch() string {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short")
	out, err := cmd.Output()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(out)), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	// Fallback: check if main exists
	if err := exec.Command("git", "rev-parse", "--verify", "main").Run(); err == nil {
		return "main"
	}
	return "master"
}

// IsInsideWorkTree returns true if the current directory is inside a git repo.
func IsInsideWorkTree() bool {
	err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run()
	return err == nil
}
