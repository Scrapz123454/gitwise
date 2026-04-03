package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider      string `yaml:"provider"`
	Model         string `yaml:"model"`
	Language      string `yaml:"language"`
	Convention    string `yaml:"convention"`
	MaxDiffLines  int    `yaml:"max_diff_lines"`
	ScopeFromPath bool   `yaml:"scope_from_path"`
	Emoji         bool   `yaml:"emoji"`
	SignCommits   bool   `yaml:"sign_commits"`

	// API keys (can also come from env vars)
	OpenAIKey    string `yaml:"openai_api_key"`
	AnthropicKey string `yaml:"anthropic_api_key"`
	GeminiKey    string `yaml:"gemini_api_key"`
	OllamaURL    string `yaml:"ollama_url"`

	// Custom OpenAI-compatible endpoint
	CustomEndpoint string `yaml:"custom_endpoint"`
	CustomAPIKey   string `yaml:"custom_api_key"`

	// Cost display
	ShowCost bool `yaml:"show_cost"`

	// Caching
	CacheEnabled bool `yaml:"cache_enabled"`
	CacheDir     string `yaml:"cache_dir"`
}

func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		Provider:      "ollama",
		Model:         "",
		Language:      "english",
		Convention:    "conventional",
		MaxDiffLines:  5000,
		ScopeFromPath: true,
		Emoji:         false,
		SignCommits:   false,
		OllamaURL:     "http://localhost:11434",
		ShowCost:      true,
		CacheEnabled:  false,
		CacheDir:      filepath.Join(home, ".cache", "gitwise"),
	}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	configPath := configFilePath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.applyEnvOverrides()
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.applyEnvOverrides()
	return cfg, nil
}

func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("GITWISE_PROVIDER"); v != "" {
		c.Provider = v
	}
	if v := os.Getenv("GITWISE_MODEL"); v != "" {
		c.Model = v
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		c.OpenAIKey = v
	}
	if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
		c.AnthropicKey = v
	}
	if v := os.Getenv("GEMINI_API_KEY"); v != "" {
		c.GeminiKey = v
	}
	if v := os.Getenv("OLLAMA_URL"); v != "" {
		c.OllamaURL = v
	}
	if v := os.Getenv("GITWISE_CUSTOM_ENDPOINT"); v != "" {
		c.CustomEndpoint = v
	}
	if v := os.Getenv("GITWISE_CUSTOM_API_KEY"); v != "" {
		c.CustomAPIKey = v
	}
}

func configFilePath() string {
	if v := os.Getenv("GITWISE_CONFIG"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gitwise", "config.yaml")
}

// LoadIgnorePatterns reads .gitwiseignore from the repo root and returns extra patterns.
func LoadIgnorePatterns(repoRoot string) []string {
	paths := []string{
		filepath.Join(repoRoot, ".gitwiseignore"),
	}

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		defer f.Close()

		var patterns []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				patterns = append(patterns, line)
			}
		}
		return patterns
	}

	return nil
}
