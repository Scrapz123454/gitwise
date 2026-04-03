package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Provider != "ollama" {
		t.Errorf("default provider = %q, want 'ollama'", cfg.Provider)
	}
	if cfg.Convention != "conventional" {
		t.Errorf("default convention = %q, want 'conventional'", cfg.Convention)
	}
	if cfg.MaxDiffLines != 5000 {
		t.Errorf("default max_diff_lines = %d, want 5000", cfg.MaxDiffLines)
	}
	if !cfg.ScopeFromPath {
		t.Error("default scope_from_path should be true")
	}
	if cfg.Emoji {
		t.Error("default emoji should be false")
	}
	if cfg.Language != "english" {
		t.Errorf("default language = %q, want 'english'", cfg.Language)
	}
}

func TestLoadIgnorePatterns(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitwise-ignore-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// No file
	patterns := LoadIgnorePatterns(dir)
	if len(patterns) != 0 {
		t.Errorf("expected 0 patterns without .gitwiseignore, got %d", len(patterns))
	}

	// With file
	content := "# comment\n*.gen.go\ngenerated/*\n\n  \n*.snap\n"
	os.WriteFile(filepath.Join(dir, ".gitwiseignore"), []byte(content), 0644)

	patterns = LoadIgnorePatterns(dir)
	expected := []string{"*.gen.go", "generated/*", "*.snap"}
	if len(patterns) != len(expected) {
		t.Fatalf("expected %d patterns, got %d: %v", len(expected), len(patterns), patterns)
	}
	for i, p := range patterns {
		if p != expected[i] {
			t.Errorf("pattern[%d] = %q, want %q", i, p, expected[i])
		}
	}
}

func TestEnvOverrides(t *testing.T) {
	cfg := DefaultConfig()

	os.Setenv("GITWISE_PROVIDER", "openai")
	os.Setenv("GITWISE_MODEL", "gpt-4o")
	defer os.Unsetenv("GITWISE_PROVIDER")
	defer os.Unsetenv("GITWISE_MODEL")

	cfg.applyEnvOverrides()

	if cfg.Provider != "openai" {
		t.Errorf("provider after env override = %q, want 'openai'", cfg.Provider)
	}
	if cfg.Model != "gpt-4o" {
		t.Errorf("model after env override = %q, want 'gpt-4o'", cfg.Model)
	}
}
