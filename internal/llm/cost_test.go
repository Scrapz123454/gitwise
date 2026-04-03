package llm

import "testing"

func TestEstimateTokens(t *testing.T) {
	// Empty string
	if tokens := EstimateTokens(""); tokens != 0 {
		t.Errorf("EstimateTokens('') = %d, want 0", tokens)
	}

	// Known text — should return > 0
	tokens := EstimateTokens("Hello, world! This is a test.")
	if tokens <= 0 {
		t.Errorf("EstimateTokens should return > 0 for non-empty text, got %d", tokens)
	}
}

func TestEstimateCost(t *testing.T) {
	// Known model
	cost := EstimateCost("gpt-4o", 1000, 200)
	if cost == "" {
		t.Error("expected non-empty cost for gpt-4o")
	}

	// Unknown model returns empty
	cost = EstimateCost("unknown-model", 1000, 200)
	if cost != "" {
		t.Errorf("expected empty cost for unknown model, got %q", cost)
	}

	// Zero tokens
	cost = EstimateCost("gpt-4o", 0, 0)
	if cost == "" {
		t.Error("expected non-empty cost string even for 0 tokens")
	}
}
