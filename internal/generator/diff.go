package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// GenerateDiffSummary generates a plain-English summary of staged changes.
func GenerateDiffSummary(cfg *config.Config) (string, error) {
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

	prompt := llm.DiffSummaryPrompt(diff)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	summary, err := provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate diff summary via %s: %w", provider.Name(), err)
	}

	return summary, nil
}
