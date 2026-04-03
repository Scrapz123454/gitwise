package git

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Workspace represents a detected monorepo workspace/package.
type Workspace struct {
	Name string
	Path string
}

// DetectWorkspaces finds monorepo workspace boundaries.
func DetectWorkspaces(repoRoot string) []Workspace {
	var workspaces []Workspace

	// Check npm/yarn/pnpm workspaces
	workspaces = append(workspaces, detectNpmWorkspaces(repoRoot)...)

	// Check Go workspaces (go.work)
	workspaces = append(workspaces, detectGoWorkspaces(repoRoot)...)

	// Check Cargo workspaces
	workspaces = append(workspaces, detectCargoWorkspaces(repoRoot)...)

	return workspaces
}

// ScopeFromWorkspace determines the workspace scope for a set of files.
func ScopeFromWorkspace(files []string, workspaces []Workspace) string {
	if len(workspaces) == 0 {
		return ""
	}

	counts := make(map[string]int)
	for _, f := range files {
		for _, ws := range workspaces {
			if strings.HasPrefix(filepath.ToSlash(f), filepath.ToSlash(ws.Path)+"/") {
				counts[ws.Name]++
			}
		}
	}

	best := ""
	bestCount := 0
	for name, count := range counts {
		if count > bestCount {
			best = name
			bestCount = count
		}
	}

	return best
}

func detectNpmWorkspaces(root string) []Workspace {
	pkgPath := filepath.Join(root, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	var pkg struct {
		Workspaces interface{} `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	var patterns []string
	switch v := pkg.Workspaces.(type) {
	case []interface{}:
		for _, p := range v {
			if s, ok := p.(string); ok {
				patterns = append(patterns, s)
			}
		}
	case map[string]interface{}:
		if pkgs, ok := v["packages"].([]interface{}); ok {
			for _, p := range pkgs {
				if s, ok := p.(string); ok {
					patterns = append(patterns, s)
				}
			}
		}
	}

	var workspaces []Workspace
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(root, pattern))
		for _, m := range matches {
			rel, _ := filepath.Rel(root, m)
			name := filepath.Base(m)
			// Try to read package name from package.json
			subPkg := filepath.Join(m, "package.json")
			if data, err := os.ReadFile(subPkg); err == nil {
				var p struct {
					Name string `json:"name"`
				}
				if json.Unmarshal(data, &p) == nil && p.Name != "" {
					name = p.Name
				}
			}
			workspaces = append(workspaces, Workspace{Name: name, Path: rel})
		}
	}

	return workspaces
}

func detectGoWorkspaces(root string) []Workspace {
	workPath := filepath.Join(root, "go.work")
	data, err := os.ReadFile(workPath)
	if err != nil {
		return nil
	}

	var workspaces []Workspace
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "use ") || strings.HasPrefix(line, "\t") {
			dir := strings.TrimSpace(strings.TrimPrefix(line, "use"))
			dir = strings.Trim(dir, "()")
			dir = strings.TrimSpace(dir)
			if dir != "" && dir != "." {
				workspaces = append(workspaces, Workspace{
					Name: filepath.Base(dir),
					Path: dir,
				})
			}
		}
	}

	return workspaces
}

func detectCargoWorkspaces(root string) []Workspace {
	cargoPath := filepath.Join(root, "Cargo.toml")
	data, err := os.ReadFile(cargoPath)
	if err != nil {
		return nil
	}

	content := string(data)
	idx := strings.Index(content, "[workspace]")
	if idx == -1 {
		return nil
	}

	var workspaces []Workspace
	membersIdx := strings.Index(content[idx:], "members")
	if membersIdx == -1 {
		return nil
	}

	// Simple TOML array parsing
	section := content[idx+membersIdx:]
	start := strings.Index(section, "[")
	end := strings.Index(section, "]")
	if start == -1 || end == -1 {
		return nil
	}

	members := section[start+1 : end]
	for _, m := range strings.Split(members, ",") {
		m = strings.TrimSpace(m)
		m = strings.Trim(m, `"'`)
		if m != "" {
			// Expand globs
			matches, _ := filepath.Glob(filepath.Join(root, m))
			if len(matches) == 0 {
				workspaces = append(workspaces, Workspace{
					Name: filepath.Base(m),
					Path: m,
				})
			}
			for _, match := range matches {
				rel, _ := filepath.Rel(root, match)
				workspaces = append(workspaces, Workspace{
					Name: filepath.Base(match),
					Path: rel,
				})
			}
		}
	}

	return workspaces
}
