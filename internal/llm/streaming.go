package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamCallback is called for each token received during streaming.
type StreamCallback func(token string)

// StreamingProvider extends Provider with streaming support.
type StreamingProvider interface {
	Provider
	GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error)
}

// --- Ollama Streaming ---

func (o *Ollama) GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: true,
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

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk ollamaResponse
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		full.WriteString(chunk.Response)
		if callback != nil {
			callback(chunk.Response)
		}
	}

	return full.String(), scanner.Err()
}

// --- OpenAI Streaming ---

type openAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (o *OpenAI) GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error) {
	return openAICompatibleStream(ctx, "https://api.openai.com/v1/chat/completions", o.apiKey, o.model, prompt, callback)
}

// --- Custom Streaming ---

func (c *Custom) GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error) {
	url := c.endpoint + "/v1/chat/completions"
	return openAICompatibleStream(ctx, url, c.apiKey, c.model, prompt, callback)
}

// openAICompatibleStream handles SSE streaming for OpenAI-compatible APIs.
func openAICompatibleStream(ctx context.Context, url, apiKey, model, prompt string, callback StreamCallback) (string, error) {
	type streamReq struct {
		Model    string          `json:"model"`
		Messages []openAIMessage `json:"messages"`
		Stream   bool            `json:"stream"`
	}

	reqBody := streamReq{
		Model:    model,
		Messages: []openAIMessage{{Role: "user", Content: prompt}},
		Stream:   true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("streaming request returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk openAIStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			token := chunk.Choices[0].Delta.Content
			full.WriteString(token)
			if callback != nil {
				callback(token)
			}
		}
	}

	return full.String(), scanner.Err()
}

// --- Anthropic Streaming ---

type anthropicStreamEvent struct {
	Type  string `json:"type"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

func (a *Anthropic) GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error) {
	reqBody := anthropicRequest{
		Model:     a.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	// Add stream field
	type streamReq struct {
		anthropicRequest
		Stream bool `json:"stream"`
	}

	sr := streamReq{anthropicRequest: reqBody, Stream: true}
	body, err := json.Marshal(sr)
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

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var event anthropicStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if event.Type == "content_block_delta" && event.Delta != nil && event.Delta.Text != "" {
			full.WriteString(event.Delta.Text)
			if callback != nil {
				callback(event.Delta.Text)
			}
		}
	}

	return full.String(), scanner.Err()
}

// --- Gemini (no SSE streaming, falls back to non-streaming) ---

func (g *Gemini) GenerateStream(ctx context.Context, prompt string, callback StreamCallback) (string, error) {
	// Gemini REST API streaming is more complex; fall back to non-streaming
	result, err := g.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}
	if callback != nil {
		callback(result)
	}
	return result, nil
}
