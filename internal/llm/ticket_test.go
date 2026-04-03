package llm

import "testing"

func TestExtractTicket(t *testing.T) {
	tests := []struct {
		branch   string
		expected string
	}{
		{"feat/SHOP-123-add-cart", "SHOP-123"},
		{"fix/JIRA-456-login-bug", "JIRA-456"},
		{"feat/PROJ-1-init", "PROJ-1"},
		{"fix/#42-bug", "#42"},
		{"feature/#100-improvement", "#100"},
		// No ticket
		{"main", ""},
		{"develop", ""},
		{"feat/add-login", ""},
		{"fix/bug-fix", ""},
		// Edge cases
		{"AB-1", "AB-1"},
		{"feat/AB-1", "AB-1"},
	}

	for _, tt := range tests {
		result := ExtractTicket(tt.branch)
		if result != tt.expected {
			t.Errorf("ExtractTicket(%q) = %q, want %q", tt.branch, result, tt.expected)
		}
	}
}
