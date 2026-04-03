package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/llm"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var amendCmd = &cobra.Command{
	Use:   "amend",
	Short: "Regenerate the message for the last commit",
	Long:  `Generates a new commit message for the last commit based on its diff, without changing the commit content.`,
	RunE:  runAmend,
}

func init() {
	rootCmd.AddCommand(amendCmd)
}

func runAmend(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	// Get the diff of the last commit
	diffCmd := exec.Command("git", "diff", "HEAD~1..HEAD")
	diffOut, err := diffCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get last commit diff: %w", err)
	}

	if len(diffOut) == 0 {
		return fmt.Errorf("no diff found for last commit")
	}

	diff := string(diffOut)

	// Get current message
	msgCmd := exec.Command("git", "log", "-1", "--format=%B")
	msgOut, err := msgCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get last commit message: %w", err)
	}
	fmt.Printf("Current message: %s\n", string(msgOut))

	provider, err := llm.NewProvider(cfg)
	if err != nil {
		return err
	}

	scope := ""
	if cfg.ScopeFromPath {
		// Get files from last commit
		filesCmd := exec.Command("git", "diff", "--name-only", "HEAD~1..HEAD")
		filesOut, _ := filesCmd.Output()
		if len(filesOut) > 0 {
			files := splitLines(string(filesOut))
			scope = git.InferScope(files)
		}
	}

	commitHistory := git.RecentCommitMessages(10)
	ticket := ""
	if branch, err := git.CurrentBranch(); err == nil {
		ticket = llm.ExtractTicket(branch)
	}

	prompt := llm.CommitPrompt(diff, cfg.Convention, cfg.ScopeFromPath, scope, cfg.Language, commitHistory, ticket)

	spinner := ui.NewSpinner("Generating new commit message...")
	spinner.Start()

	ctx, cancel := contextWithTimeout(60)
	defer cancel()

	newMsg, err := provider.Generate(ctx, prompt)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("failed to generate message: %w", err)
	}

	newMsg = cleanMsg(newMsg)
	if cfg.Emoji {
		newMsg = llm.AddEmoji(newMsg)
	}

	fmt.Printf("\nNew message: %s\n\n", newMsg)

	if !ui.Confirm("Amend with this message?") {
		fmt.Println("Aborted.")
		return nil
	}

	amendGitCmd := exec.Command("git", "commit", "--amend", "-m", newMsg)
	amendGitCmd.Stdout = os.Stdout
	amendGitCmd.Stderr = os.Stderr

	if err := amendGitCmd.Run(); err != nil {
		return fmt.Errorf("git amend failed: %w", err)
	}

	return nil
}

func splitLines(s string) []string {
	var lines []string
	for _, line := range splitString(s, "\n") {
		line = trimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
