package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Gemini struct {
	apiKey string
	model  string
}

func NewGemini(apiKey, model string) *Gemini {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &Gemini{apiKey: apiKey, model: model}
}

func (g *Gemini) Name() string {
	return fmt.Sprintf("gemini (%s)", g.model)
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (g *Gemini) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != nil {
		return "", fmt.Errorf("gemini error: %s", result.Error.Message)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("gemini returned no content")
	}
	if len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no content parts")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
