package git

import (
	"os/exec"
	"strings"
)

// BranchCommits returns the commit log for the current branch since it diverged from base.
func BranchCommits(base string) (string, error) {
	if base == "" {
		base = "main"
	}

	// Try to find the merge base
	mergeBaseCmd := exec.Command("git", "merge-base", base, "HEAD")
	mergeBaseOut, err := mergeBaseCmd.Output()
	if err != nil {
		// If merge-base fails, try with "master"
		if base == "main" {
			mergeBaseCmd = exec.Command("git", "merge-base", "master", "HEAD")
			mergeBaseOut, err = mergeBaseCmd.Output()
			if err != nil {
				return "", err
			}
			base = "master"
		} else {
			return "", err
		}
	}

	mergeBase := strings.TrimSpace(string(mergeBaseOut))

	// Get commits since merge base
	logCmd := exec.Command("git", "log", "--oneline", mergeBase+"..HEAD")
	out, err := logCmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// BranchDiff returns the combined diff of all changes on the current branch vs base.
func BranchDiff(base string, maxLines int) (string, error) {
	if base == "" {
		base = "main"
	}

	diffCmd := exec.Command("git", "diff", base+"...HEAD")
	out, err := diffCmd.Output()
	if err != nil {
		// Try master
		if base == "main" {
			diffCmd = exec.Command("git", "diff", "master...HEAD")
			out, err = diffCmd.Output()
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	diff := string(out)
	lines := strings.Split(diff, "\n")
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "\n... (diff truncated)")
		diff = strings.Join(lines, "\n")
	}

	return diff, nil
}
