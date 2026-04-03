package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PRTemplate reads the PR template from the repository if it exists.
// It checks common locations for GitHub PR templates.
func PRTemplate() string {
	root, err := repoRoot()
	if err != nil {
		return ""
	}

	// Check common PR template locations in priority order
	paths := []string{
		filepath.Join(root, ".github", "pull_request_template.md"),
		filepath.Join(root, ".github", "PULL_REQUEST_TEMPLATE.md"),
		filepath.Join(root, "pull_request_template.md"),
		filepath.Join(root, "PULL_REQUEST_TEMPLATE.md"),
		filepath.Join(root, "docs", "pull_request_template.md"),
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
