package generator

import "testing"

func TestCleanCommitMessage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cleanup
		{"  feat: add login  ", "feat: add login"},
		// Code fences
		{"```feat: add login```", "feat: add login"},
		{"```\nfeat: add login\n```", "feat: add login"},
		// Prefixes
		{"Commit message: feat: add login", "feat: add login"},
		{"commit: feat: add login", "feat: add login"},
		{"Message: fix: resolve crash", "fix: resolve crash"},
		// LLM preamble
		{"Here is the commit message:\nfeat: add login", "feat: add login"},
		{"Here's the commit message:\nfix: bug", "fix: bug"},
		{"Based on the diff:\nrefactor: simplify", "refactor: simplify"},
		// No cleanup needed
		{"feat(auth): add JWT", "feat(auth): add JWT"},
		// Empty
		{"", ""},
		{"   ", ""},
	}

	for _, tt := range tests {
		result := cleanCommitMessage(tt.input)
		if result != tt.expected {
			t.Errorf("cleanCommitMessage(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
