package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// RecentCommitMessages returns the last N commit messages for style matching.
func RecentCommitMessages(count int) string {
	cmd := exec.Command("git", "log", "--oneline", "--no-merges", "-n", fmt.Sprintf("%d", count))
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CommitsBetweenTags returns commit log between two tags.
func CommitsBetweenTags(fromTag, toTag string) (string, error) {
	var args []string
	if fromTag != "" {
		args = []string{"log", "--oneline", "--no-merges", fromTag + ".." + toTag}
	} else {
		args = []string{"log", "--oneline", "--no-merges", toTag}
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// LatestTag returns the most recent tag.
func LatestTag() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// AllTags returns all tags sorted by version.
func AllTags() []string {
	cmd := exec.Command("git", "tag", "--sort=-v:refname")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil
	}
	return strings.Split(raw, "\n")
}

// CommitMessagesForRange returns detailed commit messages for linting.
func CommitMessagesForRange(commitRange string) (string, error) {
	args := []string{"log", "--format=%H %s", "--no-merges"}
	if commitRange != "" {
		args = append(args, commitRange)
	} else {
		args = append(args, "-n", "10")
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
