package git

import "testing"

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		filename string
		ignored  bool
	}{
		{"package-lock.json", true},
		{"yarn.lock", true},
		{"go.sum", true},
		{"Cargo.lock", true},
		{"app.min.js", true},
		{"style.min.css", true},
		{"bundle.js.map", true},
		{"model.pb.go", true},
		{"vendor/lib/x.go", true},
		{"node_modules/foo/index.js", true},
		{"dist/bundle.js", true},
		// Not ignored
		{"main.go", false},
		{"README.md", false},
		{"src/app.ts", false},
		{"internal/config/config.go", false},
		{"package.json", false},
	}

	for _, tt := range tests {
		result := shouldIgnore(tt.filename)
		if result != tt.ignored {
			t.Errorf("shouldIgnore(%q) = %v, want %v", tt.filename, result, tt.ignored)
		}
	}
}

func TestInferScope(t *testing.T) {
	tests := []struct {
		files    []string
		expected string
	}{
		{[]string{"internal/auth/login.go", "internal/auth/jwt.go"}, "internal"},
		{[]string{"cmd/commit.go"}, "cmd"},
		{[]string{"README.md"}, ""},
		{nil, ""},
		{[]string{}, ""},
	}

	for _, tt := range tests {
		result := InferScope(tt.files)
		if result != tt.expected {
			t.Errorf("InferScope(%v) = %q, want %q", tt.files, result, tt.expected)
		}
	}
}

func TestSetExtraIgnorePatterns(t *testing.T) {
	// Reset state
	SetExtraIgnorePatterns(nil)

	if shouldIgnore("custom.gen.go") {
		t.Error("custom.gen.go should not be ignored before adding pattern")
	}

	SetExtraIgnorePatterns([]string{"*.gen.go"})

	if !shouldIgnore("custom.gen.go") {
		t.Error("custom.gen.go should be ignored after adding pattern")
	}

	// Cleanup
	SetExtraIgnorePatterns(nil)
}
