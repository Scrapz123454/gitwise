package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// GenerateReview performs AI code review on staged changes.
func GenerateReview(cfg *config.Config) (string, error) {
	diff, err := git.StagedDiff(cfg.MaxDiffLines)
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	if diff == "" {
		return "", fmt.Errorf("no staged changes found — stage your changes with `git add` first")
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return "", err
	}

	prompt := llm.ReviewPrompt(diff)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	review, err := provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate review via %s: %w", provider.Name(), err)
	}

	return review, nil
}
