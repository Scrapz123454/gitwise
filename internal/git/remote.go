package git

import (
	"os/exec"
	"strings"
)

// RemoteURL returns the remote URL for the given remote name.
func RemoteURL(name string) (string, error) {
	if name == "" {
		name = "origin"
	}
	cmd := exec.Command("git", "remote", "get-url", name)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
