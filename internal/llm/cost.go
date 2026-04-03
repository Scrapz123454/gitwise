package llm

import (
	"fmt"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

// TokenCost holds pricing per 1M tokens for a model.
type TokenCost struct {
	InputPerM  float64
	OutputPerM float64
}

// Known pricing (approximate, per 1M tokens).
var modelPricing = map[string]TokenCost{
	// OpenAI
	"gpt-4o":          {InputPerM: 2.50, OutputPerM: 10.0},
	"gpt-4o-mini":     {InputPerM: 0.15, OutputPerM: 0.60},
	"gpt-4-turbo":     {InputPerM: 10.0, OutputPerM: 30.0},
	// Anthropic
	"claude-sonnet-4-20250514":  {InputPerM: 3.0, OutputPerM: 15.0},
	"claude-haiku-4-5-20251001": {InputPerM: 0.80, OutputPerM: 4.0},
	// Gemini
	"gemini-2.0-flash":   {InputPerM: 0.10, OutputPerM: 0.40},
	"gemini-2.5-pro":     {InputPerM: 1.25, OutputPerM: 10.0},
}

// EstimateTokens counts approximate tokens in text using tiktoken.
func EstimateTokens(text string) int {
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		// Fallback: rough estimate of 1 token per 4 chars
		return len(text) / 4
	}
	return len(enc.Encode(text, nil, nil))
}

// EstimateCost returns the estimated cost for a prompt given the model.
func EstimateCost(model string, inputTokens, outputTokens int) string {
	pricing, ok := modelPricing[model]
	if !ok {
		return ""
	}

	inputCost := float64(inputTokens) / 1_000_000 * pricing.InputPerM
	outputCost := float64(outputTokens) / 1_000_000 * pricing.OutputPerM
	total := inputCost + outputCost

	if total < 0.001 {
		return fmt.Sprintf("~$0.00 (%d input tokens)", inputTokens)
	}
	return fmt.Sprintf("~$%.4f (%d input tokens)", total, inputTokens)
}
