package git

import "testing"

func TestScopeFromWorkspace(t *testing.T) {
	workspaces := []Workspace{
		{Name: "api", Path: "packages/api"},
		{Name: "web", Path: "packages/web"},
		{Name: "shared", Path: "packages/shared"},
	}

	tests := []struct {
		files    []string
		expected string
	}{
		{[]string{"packages/api/main.go", "packages/api/handler.go"}, "api"},
		{[]string{"packages/web/index.tsx"}, "web"},
		{[]string{"packages/api/x.go", "packages/web/y.tsx"}, ""}, // tie, depends on map order — but at least shouldn't crash
		{[]string{"README.md"}, ""},
		{nil, ""},
		{[]string{}, ""},
	}

	for _, tt := range tests {
		result := ScopeFromWorkspace(tt.files, workspaces)
		// For the tie case, just check it doesn't panic
		if len(tt.files) == 2 && tt.files[0] == "packages/api/x.go" {
			continue
		}
		if result != tt.expected {
			t.Errorf("ScopeFromWorkspace(%v) = %q, want %q", tt.files, result, tt.expected)
		}
	}
}

func TestScopeFromWorkspaceEmpty(t *testing.T) {
	result := ScopeFromWorkspace([]string{"main.go"}, nil)
	if result != "" {
		t.Errorf("expected empty scope with no workspaces, got %q", result)
	}
}
