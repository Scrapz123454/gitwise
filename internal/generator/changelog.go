package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// GenerateChangelog generates a CHANGELOG entry from commits between tags.
func GenerateChangelog(cfg *config.Config, fromTag, toTag string) (string, error) {
	if toTag == "" {
		toTag = "HEAD"
	}

	commits, err := git.CommitsBetweenTags(fromTag, toTag)
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}
	if commits == "" {
		return "", fmt.Errorf("no commits found between %s and %s", fromTag, toTag)
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return "", err
	}

	prompt := llm.ChangelogPrompt(commits, fromTag, toTag)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	changelog, err := provider.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate changelog via %s: %w", provider.Name(), err)
	}

	return changelog, nil
}
