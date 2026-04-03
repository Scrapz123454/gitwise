package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// LintCommits analyzes commit messages for conventional commit compliance.
func LintCommits(cfg *config.Config, commitRange string) (string, error) {
	commits, err := git.CommitMessagesForRange(commitRange)
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}
	if commits == "" {
		return "", fmt.Errorf("no commits found in range %s", commitRange)
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return "", err
	}

	prompt := llm.LintPrompt(commits)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to lint commits via %s: %w", provider.Name(), err)
	}

	return result, nil
}
