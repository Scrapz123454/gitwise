package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard",
	Long:  `Guided setup: detect LLMs, test connections, create config, and install git hooks.`,
	RunE:  runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println("Welcome to gitwise setup!")
	fmt.Println("========================")
	fmt.Println()

	cfg := config.DefaultConfig()

	// Step 1: Detect available providers
	fmt.Println("Step 1: Detecting available LLM providers...")
	fmt.Println()

	providers := detectProviders()
	if len(providers) == 0 {
		fmt.Println("  No providers detected. You'll need to configure one manually.")
		fmt.Println()
	} else {
		for _, p := range providers {
			fmt.Printf("  [OK] %s\n", p)
		}
		fmt.Println()

		choice := ui.PromptChoice("Select your preferred provider:", providers)
		selectedProvider := providers[choice]

		switch {
		case selectedProvider == "ollama (local)" || selectedProvider == "ollama":
			cfg.Provider = "ollama"
			ollamaURL := os.Getenv("OLLAMA_URL")
			if ollamaURL == "" {
				ollamaURL = "http://localhost:11434"
			}
			models, err := llm.ListOllamaModels(ollamaURL)
			if err == nil && len(models) > 0 {
				fmt.Printf("  Found %d models:\n", len(models))
				modelChoice := ui.PromptChoice("Select a model:", models)
				cfg.Model = models[modelChoice]
			} else {
				fmt.Print("  Model name [llama3]: ")
				var model string
				fmt.Scanln(&model)
				if model != "" {
					cfg.Model = model
				} else {
					cfg.Model = "llama3"
				}
			}
		case selectedProvider == "openai (OPENAI_API_KEY found)":
			cfg.Provider = "openai"
			cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
			cfg.Model = "gpt-4o-mini"
		case selectedProvider == "anthropic (ANTHROPIC_API_KEY found)":
			cfg.Provider = "anthropic"
			cfg.AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")
			cfg.Model = "claude-sonnet-4-20250514"
		case selectedProvider == "gemini (GEMINI_API_KEY found)":
			cfg.Provider = "gemini"
			cfg.GeminiKey = os.Getenv("GEMINI_API_KEY")
			cfg.Model = "gemini-2.0-flash"
		}
	}

	// Step 2: Convention
	fmt.Println()
	fmt.Println("Step 2: Commit message convention")
	conventionChoice := ui.PromptChoice("Choose commit convention:", []string{
		"conventional (feat/fix/refactor...)",
		"angular (same as conventional)",
		"none (free-form)",
	})
	switch conventionChoice {
	case 0:
		cfg.Convention = "conventional"
	case 1:
		cfg.Convention = "angular"
	case 2:
		cfg.Convention = "none"
	}

	// Step 3: Emoji
	fmt.Println()
	cfg.Emoji = ui.Confirm("Step 3: Use gitmoji (emoji prefixes)?")

	// Step 4: Language
	fmt.Println()
	fmt.Print("Step 4: Commit message language [english]: ")
	var lang string
	fmt.Scanln(&lang)
	if lang != "" {
		cfg.Language = lang
	}

	// Step 5: Save config
	fmt.Println()
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "gitwise")
	configPath := filepath.Join(configDir, "config.yaml")

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("  Config saved to %s\n", configPath)

	// Step 6: Install hook
	fmt.Println()
	if ui.Confirm("Step 5: Install git hook (auto-generate on `git commit`)?") {
		if err := hookInstallCmd.RunE(cmd, nil); err != nil {
			fmt.Printf("  Warning: hook install failed: %s\n", err)
		}
	}

	fmt.Println()
	fmt.Println("Setup complete! Try `gitwise commit` to generate your first message.")
	fmt.Println()

	return nil
}

func detectProviders() []string {
	var providers []string

	// Check Ollama
	client := &http.Client{Timeout: 2 * time.Second}
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	resp, err := client.Get(ollamaURL + "/api/tags")
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode == 200 {
			providers = append(providers, "ollama (local)")
		}
	}

	// Check API keys
	if os.Getenv("OPENAI_API_KEY") != "" {
		providers = append(providers, "openai (OPENAI_API_KEY found)")
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		providers = append(providers, "anthropic (ANTHROPIC_API_KEY found)")
	}
	if os.Getenv("GEMINI_API_KEY") != "" {
		providers = append(providers, "gemini (GEMINI_API_KEY found)")
	}

	// Check custom endpoint
	if os.Getenv("GITWISE_CUSTOM_ENDPOINT") != "" {
		providers = append(providers, "custom (GITWISE_CUSTOM_ENDPOINT found)")
	}

	// Check LM Studio (common port)
	resp2, err := client.Get("http://localhost:1234/v1/models")
	if err == nil {
		resp2.Body.Close()
		if resp2.StatusCode == 200 {
			providers = append(providers, "custom (LM Studio detected at :1234)")
		}
	}

	return providers
}

