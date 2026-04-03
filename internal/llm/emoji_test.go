package llm

import "testing"

func TestAddEmoji(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"feat(auth): add login", "✨ feat(auth): add login"},
		{"feat: add login", "✨ feat: add login"},
		{"fix: resolve crash", "🐛 fix: resolve crash"},
		{"fix(api): null pointer", "🐛 fix(api): null pointer"},
		{"docs: update README", "📝 docs: update README"},
		{"refactor(core): simplify", "♻️ refactor(core): simplify"},
		{"test: add unit tests", "✅ test: add unit tests"},
		{"perf: optimize query", "⚡ perf: optimize query"},
		{"chore: update deps", "🔧 chore: update deps"},
		{"ci: add workflow", "👷 ci: add workflow"},
		{"build: update makefile", "📦 build: update makefile"},
		{"style: format code", "💄 style: format code"},
		{"revert: undo change", "⏪ revert: undo change"},
		// No emoji for unknown types
		{"plain message", "plain message"},
		{"unknown: something", "unknown: something"},
		{"", ""},
	}

	for _, tt := range tests {
		result := AddEmoji(tt.input)
		if result != tt.expected {
			t.Errorf("AddEmoji(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
