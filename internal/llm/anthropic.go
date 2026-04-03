package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Anthropic struct {
	apiKey string
	model  string
}

func NewAnthropic(apiKey, model string) *Anthropic {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &Anthropic{apiKey: apiKey, model: model}
}

func (a *Anthropic) Name() string {
	return fmt.Sprintf("anthropic (%s)", a.model)
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *Anthropic) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := anthropicRequest{
		Model:     a.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Anthropic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("anthropic error: %s", result.Error.Message)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic returned no content")
	}

	return result.Content[0].Text, nil
}
