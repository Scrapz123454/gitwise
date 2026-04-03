package git

import (
	"os/exec"
	"strings"
)

// StatusShort returns a short git status.
func StatusShort() (string, error) {
	cmd := exec.Command("git", "status", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
