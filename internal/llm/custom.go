package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Custom implements an OpenAI-compatible API client for any endpoint
// (Groq, Together, LM Studio, vLLM, etc.)
type Custom struct {
	endpoint string
	apiKey   string
	model    string
}

func NewCustom(endpoint, apiKey, model string) *Custom {
	endpoint = strings.TrimSuffix(endpoint, "/")
	return &Custom{endpoint: endpoint, apiKey: apiKey, model: model}
}

func (c *Custom) Name() string {
	return fmt.Sprintf("custom (%s @ %s)", c.model, c.endpoint)
}

func (c *Custom) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := openAIRequest{
		Model: c.model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := c.endpoint + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("custom endpoint returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("custom endpoint error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("custom endpoint returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}
