package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const hookScript = `#!/bin/sh
# gitwise prepare-commit-msg hook
# Generates AI-powered commit messages

# Only run if no message was provided (not amend, merge, etc.)
if [ "$2" = "" ]; then
    MSG=$(gitwise commit --hook 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$MSG" ]; then
        echo "$MSG" > "$1"
    fi
fi
`

// InstallHook installs gitwise as a prepare-commit-msg hook.
func InstallHook() error {
	hooksDir, err := gitHooksDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	// Check if hook already exists
	if data, err := os.ReadFile(hookPath); err == nil {
		if strings.Contains(string(data), "gitwise") {
			return fmt.Errorf("gitwise hook is already installed")
		}
		return fmt.Errorf("a prepare-commit-msg hook already exists — use --force to overwrite")
	}

	return os.WriteFile(hookPath, []byte(hookScript), 0o755)
}

// UninstallHook removes the gitwise prepare-commit-msg hook.
func UninstallHook() error {
	hooksDir, err := gitHooksDir()
	if err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")
	data, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("no prepare-commit-msg hook found")
	}

	if !strings.Contains(string(data), "gitwise") {
		return fmt.Errorf("the existing hook was not installed by gitwise")
	}

	return os.Remove(hookPath)
}

func gitHooksDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git repository")
	}
	return filepath.Join(strings.TrimSpace(string(out)), "hooks"), nil
}
