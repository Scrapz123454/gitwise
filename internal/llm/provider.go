package llm

import (
	"context"
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
)

// Provider is the interface that all LLM backends must implement.
type Provider interface {
	// Generate sends a prompt and returns the generated text.
	Generate(ctx context.Context, prompt string) (string, error)

	// Name returns the provider name for display.
	Name() string
}

// NewProvider creates an LLM provider based on configuration.
func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "ollama":
		return NewOllama(cfg.OllamaURL, cfg.Model), nil
	case "openai":
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for the openai provider")
		}
		return NewOpenAI(cfg.OpenAIKey, cfg.Model), nil
	case "anthropic":
		if cfg.AnthropicKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY is required for the anthropic provider")
		}
		return NewAnthropic(cfg.AnthropicKey, cfg.Model), nil
	case "gemini":
		if cfg.GeminiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY is required for the gemini provider")
		}
		return NewGemini(cfg.GeminiKey, cfg.Model), nil
	case "custom":
		if cfg.CustomEndpoint == "" {
			return nil, fmt.Errorf("custom_endpoint is required for the custom provider (set GITWISE_CUSTOM_ENDPOINT)")
		}
		return NewCustom(cfg.CustomEndpoint, cfg.CustomAPIKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s (supported: ollama, openai, anthropic, gemini, custom)", cfg.Provider)
	}
}
