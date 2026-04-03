package generator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
)

// CommitResult holds the generated message plus metadata.
type CommitResult struct {
	Message     string
	TokenCount  int
	CostEstimate string
}

// GenerateCommitMessage generates a commit message from staged changes.
func GenerateCommitMessage(cfg *config.Config) (*CommitResult, error) {
	diff, err := git.StagedDiff(cfg.MaxDiffLines)
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}
	if diff == "" {
		return nil, fmt.Errorf("no staged changes found — stage your changes with `git add` first")
	}

	// Load .gitwiseignore
	if root, err := repoRoot(); err == nil {
		patterns := config.LoadIgnorePatterns(root)
		git.SetExtraIgnorePatterns(patterns)
	}

	scope := inferScope(cfg)
	commitHistory := ""
	if history := git.RecentCommitMessages(10); history != "" {
		commitHistory = history
	}

	// Detect ticket from branch name
	ticket := ""
	if branch, err := git.CurrentBranch(); err == nil {
		ticket = llm.ExtractTicket(branch)
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	prompt := llm.CommitPrompt(diff, cfg.Convention, cfg.ScopeFromPath, scope, cfg.Language, commitHistory, ticket)
	inputTokens := llm.EstimateTokens(prompt)

	// Check cache
	if cfg.CacheEnabled {
		cache := llm.NewCache(cfg.CacheDir)
		if cached, ok := cache.Get(prompt); ok {
			msg := cleanCommitMessage(cached)
			if cfg.Emoji {
				msg = llm.AddEmoji(msg)
			}
			return &CommitResult{
				Message:      msg,
				TokenCount:   inputTokens,
				CostEstimate: "(cached)",
			}, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	msg, err := provider.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message via %s: %w", provider.Name(), err)
	}

	// Cache the response
	if cfg.CacheEnabled {
		cache := llm.NewCache(cfg.CacheDir)
		_ = cache.Set(prompt, msg)
	}

	msg = cleanCommitMessage(msg)
	if cfg.Emoji {
		msg = llm.AddEmoji(msg)
	}

	outputTokens := llm.EstimateTokens(msg)
	costEstimate := llm.EstimateCost(cfg.Model, inputTokens, outputTokens)

	return &CommitResult{
		Message:      msg,
		TokenCount:   inputTokens,
		CostEstimate: costEstimate,
	}, nil
}

// GenerateCommitMessageStreaming generates with token-by-token streaming.
func GenerateCommitMessageStreaming(cfg *config.Config, callback llm.StreamCallback) (*CommitResult, error) {
	diff, err := git.StagedDiff(cfg.MaxDiffLines)
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}
	if diff == "" {
		return nil, fmt.Errorf("no staged changes found — stage your changes with `git add` first")
	}

	scope := inferScope(cfg)
	commitHistory := git.RecentCommitMessages(10)
	ticket := ""
	if branch, err := git.CurrentBranch(); err == nil {
		ticket = llm.ExtractTicket(branch)
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	prompt := llm.CommitPrompt(diff, cfg.Convention, cfg.ScopeFromPath, scope, cfg.Language, commitHistory, ticket)
	inputTokens := llm.EstimateTokens(prompt)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	streamer, ok := provider.(llm.StreamingProvider)
	if !ok {
		// Fallback to non-streaming
		msg, err := provider.Generate(ctx, prompt)
		if err != nil {
			return nil, err
		}
		msg = cleanCommitMessage(msg)
		if cfg.Emoji {
			msg = llm.AddEmoji(msg)
		}
		if callback != nil {
			callback(msg)
		}
		return &CommitResult{Message: msg, TokenCount: inputTokens}, nil
	}

	msg, err := streamer.GenerateStream(ctx, prompt, callback)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message via %s: %w", provider.Name(), err)
	}

	msg = cleanCommitMessage(msg)
	if cfg.Emoji {
		msg = llm.AddEmoji(msg)
	}

	outputTokens := llm.EstimateTokens(msg)
	costEstimate := llm.EstimateCost(cfg.Model, inputTokens, outputTokens)

	return &CommitResult{
		Message:      msg,
		TokenCount:   inputTokens,
		CostEstimate: costEstimate,
	}, nil
}

// GenerateMultipleCommitMessages generates multiple commit message options.
func GenerateMultipleCommitMessages(cfg *config.Config, count int) ([]string, error) {
	diff, err := git.StagedDiff(cfg.MaxDiffLines)
	if err != nil {
		return nil, fmt.Errorf("failed to get staged diff: %w", err)
	}
	if diff == "" {
		return nil, fmt.Errorf("no staged changes found — stage your changes with `git add` first")
	}

	scope := inferScope(cfg)
	commitHistory := git.RecentCommitMessages(10)
	ticket := ""
	if branch, err := git.CurrentBranch(); err == nil {
		ticket = llm.ExtractTicket(branch)
	}

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	prompt := llm.MultiCommitPrompt(diff, cfg.Convention, cfg.ScopeFromPath, scope, count, cfg.Language, commitHistory, ticket)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	raw, err := provider.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit messages via %s: %w", provider.Name(), err)
	}

	parts := strings.Split(raw, "---")
	var messages []string
	for _, p := range parts {
		msg := cleanCommitMessage(p)
		if msg != "" {
			if cfg.Emoji {
				msg = llm.AddEmoji(msg)
			}
			messages = append(messages, msg)
		}
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("failed to parse commit message options")
	}

	return messages, nil
}

func inferScope(cfg *config.Config) string {
	if !cfg.ScopeFromPath {
		return ""
	}
	files, _ := git.StagedFiles()

	// Try monorepo workspace scope first
	if root, err := repoRoot(); err == nil {
		workspaces := git.DetectWorkspaces(root)
		if ws := git.ScopeFromWorkspace(files, workspaces); ws != "" {
			return ws
		}
	}

	return git.InferScope(files)
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func cleanCommitMessage(msg string) string {
	msg = strings.TrimSpace(msg)

	// Remove markdown code fences if present
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	// Remove common LLM preamble lines before the actual commit message
	lines := strings.SplitN(msg, "\n", 2)
	if len(lines) >= 2 {
		firstLower := strings.ToLower(lines[0])
		preambles := []string{
			"here is the commit message",
			"here's the commit message",
			"here is a commit message",
			"here's a commit message",
			"sure, here",
			"based on the diff",
			"the commit message",
		}
		for _, p := range preambles {
			if strings.Contains(firstLower, p) {
				msg = strings.TrimSpace(lines[1])
				break
			}
		}
	}

	lower := strings.ToLower(msg)
	for _, prefix := range []string{"commit message:", "commit:", "message:"} {
		if strings.HasPrefix(lower, prefix) {
			msg = strings.TrimSpace(msg[len(prefix):])
			lower = strings.ToLower(msg)
		}
	}

	return msg
}
