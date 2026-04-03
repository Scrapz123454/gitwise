package generator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// PRResult holds the generated PR description plus metadata.
type PRResult struct {
	Description  string
	Labels       []string
	TokenCount   int
	CostEstimate string
}

// GeneratePRDescription generates a PR description for the current branch.
func GeneratePRDescription(cfg *config.Config, base string) (*PRResult, error) {
	if base == "" {
		base = git.DefaultBranch()
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	if branch == base {
		return nil, fmt.Errorf("you are on the base branch '%s' — switch to a feature branch first", base)
	}

	commits, err := git.BranchCommits(base)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch commits: %w", err)
	}
	if commits == "" {
		return nil, fmt.Errorf("no commits found on this branch compared to '%s'", base)
	}

	diff, err := git.BranchDiff(base, cfg.MaxDiffLines)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch diff: %w", err)
	}

	template := git.PRTemplate()

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	prompt := llm.PRPrompt(branch, base, commits, diff, template)
	inputTokens := llm.EstimateTokens(prompt)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	desc, err := provider.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PR description via %s: %w", provider.Name(), err)
	}

	outputTokens := llm.EstimateTokens(desc)
	costEstimate := llm.EstimateCost(cfg.Model, inputTokens, outputTokens)

	return &PRResult{
		Description:  desc,
		TokenCount:   inputTokens,
		CostEstimate: costEstimate,
	}, nil
}

// SuggestLabels generates label suggestions for a PR.
func SuggestLabels(cfg *config.Config, base string) ([]string, error) {
	if base == "" {
		base = git.DefaultBranch()
	}

	diff, err := git.BranchDiff(base, 500)
	if err != nil {
		return nil, err
	}

	// Get changed files from diff headers
	var files []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+++ b/") {
			files = append(files, strings.TrimPrefix(line, "+++ b/"))
		}
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	prompt := llm.LabelPrompt(files, diff)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := provider.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var labels []string
	for _, line := range strings.Split(result, "\n") {
		label := strings.TrimSpace(strings.ToLower(line))
		label = strings.TrimPrefix(label, "- ")
		label = strings.TrimPrefix(label, "* ")
		if label != "" && !strings.Contains(label, " ") {
			labels = append(labels, label)
		}
	}

	return labels, nil
}
