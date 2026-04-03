package git

import (
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// defaultIgnoredFiles are the built-in patterns for files excluded from diff analysis.
var defaultIgnoredFiles = []string{
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"go.sum",
	"Cargo.lock",
	"Gemfile.lock",
	"poetry.lock",
	"composer.lock",
	"*.min.js",
	"*.min.css",
	"*.map",
	"*.pb.go",
	"*.generated.*",
	"dist/*",
	"build/*",
	"vendor/*",
	"node_modules/*",
}

var (
	extraPatterns   []string
	extraPatternsMu sync.Mutex
)

// SetExtraIgnorePatterns adds patterns from .gitwiseignore.
func SetExtraIgnorePatterns(patterns []string) {
	extraPatternsMu.Lock()
	defer extraPatternsMu.Unlock()
	extraPatterns = patterns
}

func allIgnorePatterns() []string {
	extraPatternsMu.Lock()
	defer extraPatternsMu.Unlock()
	result := make([]string, 0, len(defaultIgnoredFiles)+len(extraPatterns))
	result = append(result, defaultIgnoredFiles...)
	result = append(result, extraPatterns...)
	return result
}

// StagedDiff returns the staged diff, filtering out ignored files.
func StagedDiff(maxLines int) (string, error) {
	filesCmd := exec.Command("git", "diff", "--staged", "--name-only")
	filesOut, err := filesCmd.Output()
	if err != nil {
		return "", err
	}

	files := strings.Split(strings.TrimSpace(string(filesOut)), "\n")
	if len(files) == 0 || (len(files) == 1 && files[0] == "") {
		return "", nil
	}

	var included []string
	for _, f := range files {
		if !shouldIgnore(f) {
			included = append(included, f)
		}
	}

	if len(included) == 0 {
		return "", nil
	}

	args := []string{"diff", "--staged", "--"}
	args = append(args, included...)
	diffCmd := exec.Command("git", args...)
	out, err := diffCmd.Output()
	if err != nil {
		return "", err
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

// StagedFiles returns the list of staged file paths.
func StagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, nil
	}
	return strings.Split(raw, "\n"), nil
}

func shouldIgnore(filename string) bool {
	base := filepath.Base(filename)
	normalized := filepath.ToSlash(filename)
	for _, pattern := range allIgnorePatterns() {
		// Match against base name (e.g., "package-lock.json")
		if matched, _ := filepath.Match(pattern, base); matched {
			return true
		}
		// Match against full path (e.g., "dist/bundle.js")
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}
		// For directory patterns like "vendor/*", check if file is under that directory
		if strings.HasSuffix(pattern, "/*") {
			dir := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(normalized, dir+"/") {
				return true
			}
		}
	}
	return false
}

// InferScope tries to detect a scope from the list of changed file paths.
func InferScope(files []string) string {
	if len(files) == 0 {
		return ""
	}

	dirs := make(map[string]int)
	for _, f := range files {
		dir := filepath.Dir(f)
		parts := strings.Split(filepath.ToSlash(dir), "/")
		for _, p := range parts {
			if p != "." && p != "" {
				dirs[p]++
				break
			}
		}
	}

	bestDir := ""
	bestCount := 0
	for dir, count := range dirs {
		if count > bestCount {
			bestDir = dir
			bestCount = count
		}
	}

	return bestDir
}

// HasStagedChanges returns true if there are any staged changes.
func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--staged", "--quiet")
	err := cmd.Run()
	return err != nil // exit code 1 means there are changes
}
