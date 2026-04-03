package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Ollama struct {
	baseURL string
	model   string
}

// visionFamilies are model families that are vision-only and not good for text generation.
var visionFamilies = map[string]bool{
	"clip":   true,
	"mllama": true,
}

// preferredModels is the priority order for auto-selecting a text model.
var preferredModels = []string{
	"llama3",
	"deepseek",
	"codellama",
	"mistral",
	"qwen",
	"gemma",
	"phi",
}

func NewOllama(baseURL, model string) *Ollama {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = detectBestOllamaModel(baseURL)
	}
	return &Ollama{baseURL: baseURL, model: model}
}

func (o *Ollama) Name() string {
	return fmt.Sprintf("ollama (%s)", o.model)
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

type ollamaTagsResponse struct {
	Models []ollamaModel `json:"models"`
}

type ollamaModel struct {
	Name    string `json:"name"`
	Details struct {
		Family   string   `json:"family"`
		Families []string `json:"families"`
	} `json:"details"`
}

// detectBestOllamaModel queries Ollama for available models and picks the best text model.
func detectBestOllamaModel(baseURL string) string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return "llama3" // fallback
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "llama3"
	}

	var tags ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "llama3"
	}

	if len(tags.Models) == 0 {
		return "llama3"
	}

	// Filter to text-only models (exclude vision-only models)
	var textModels []ollamaModel
	for _, m := range tags.Models {
		if isTextModel(m) {
			textModels = append(textModels, m)
		}
	}

	// If no text models found, try all models
	if len(textModels) == 0 {
		textModels = tags.Models
	}

	if len(textModels) == 0 {
		return "llama3"
	}

	// Pick by preference order
	for _, pref := range preferredModels {
		for _, m := range textModels {
			if strings.Contains(strings.ToLower(m.Name), pref) {
				return m.Name
			}
		}
	}

	// No preferred match — return the first text model
	return textModels[0].Name
}

// isTextModel returns true if the model is suitable for text generation.
func isTextModel(m ollamaModel) bool {
	families := m.Details.Families
	if len(families) == 0 {
		return true // assume text if no family info
	}

	// A model is vision-only if ALL its families are vision families
	hasTextFamily := false
	for _, f := range families {
		if !visionFamilies[f] {
			hasTextFamily = true
			break
		}
	}

	return hasTextFamily
}

// ListModels returns all available Ollama models.
func ListOllamaModels(baseURL string) ([]string, error) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama at %s: %w", baseURL, err)
	}
	defer resp.Body.Close()

	var tags ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	var names []string
	for _, m := range tags.Models {
		names = append(names, m.Name)
	}
	return names, nil
}

func (o *Ollama) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama at %s: %w (is Ollama running?)", o.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
