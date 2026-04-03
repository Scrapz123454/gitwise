package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message from staged changes",
	Long:  `Analyzes your staged git diff and generates a conventional commit message using AI.`,
	RunE:  runCommit,
}

func init() {
	commitCmd.Flags().BoolP("interactive", "i", false, "Review and edit the message before committing")
	commitCmd.Flags().Bool("hook", false, "Output message only (for git hook integration)")
	commitCmd.Flags().Bool("dry-run", false, "Print the message without committing")
	commitCmd.Flags().IntP("suggestions", "n", 1, "Number of commit message suggestions to generate")
	commitCmd.Flags().Bool("stream", false, "Stream output token by token")
	commitCmd.Flags().Bool("tui", false, "Use TUI picker for multiple suggestions")
	rootCmd.AddCommand(commitCmd)
}

func runCommit(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	hookMode, _ := cmd.Flags().GetBool("hook")
	interactive, _ := cmd.Flags().GetBool("interactive")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	suggestions, _ := cmd.Flags().GetInt("suggestions")
	useTUI, _ := cmd.Flags().GetBool("tui")

	var msg string

	if suggestions > 1 {
		spinner := ui.NewSpinner("Generating commit messages...")
		if !hookMode {
			spinner.Start()
		}

		messages, err := generator.GenerateMultipleCommitMessages(cfg, suggestions)
		if !hookMode {
			spinner.Stop()
		}
		if err != nil {
			return err
		}

		if hookMode {
			fmt.Print(messages[0])
			return nil
		}

		if useTUI {
			choice, err := ui.RunCommitPicker(messages)
			if err != nil {
				return err
			}
			if choice == -1 {
				fmt.Println("Aborted.")
				return nil
			}
			msg = messages[choice]
		} else {
			fmt.Println()
			choice := ui.PromptChoice("Pick a commit message:", messages)
			msg = messages[choice]
		}
	} else {
		spinner := ui.NewSpinner("Generating commit message...")
		if !hookMode {
			spinner.Start()
		}

		result, err := generator.GenerateCommitMessage(cfg)
		if !hookMode {
			spinner.Stop()
		}
		if err != nil {
			return err
		}

		msg = result.Message

		if hookMode {
			fmt.Print(msg)
			return nil
		}

		fmt.Printf("\n  %s\n", msg)
		if cfg.ShowCost && result.CostEstimate != "" {
			fmt.Printf("  %s\n", ui.FormatCostInfo(result.TokenCount, result.CostEstimate))
		}
		fmt.Println()
	}

	if dryRun {
		return nil
	}

	if interactive {
		choice := ui.PromptChoice("What would you like to do?", []string{
			"Accept and commit",
			"Edit in $EDITOR",
			"Regenerate",
			"Cancel",
		})

		switch choice {
		case 0:
			// Accept
		case 1:
			edited, err := ui.EditInEditor(msg)
			if err != nil {
				return fmt.Errorf("editor failed: %w", err)
			}
			if edited == "" {
				fmt.Println("Empty message, aborting.")
				return nil
			}
			msg = edited
		case 2:
			spinner := ui.NewSpinner("Regenerating...")
			spinner.Start()
			result, err := generator.GenerateCommitMessage(cfg)
			spinner.Stop()
			if err != nil {
				return err
			}
			msg = result.Message
			fmt.Printf("\n  %s\n\n", msg)
			if !ui.Confirm("Accept this message?") {
				fmt.Println("Aborted.")
				return nil
			}
		case 3:
			fmt.Println("Aborted.")
			return nil
		}
	} else {
		if !ui.Confirm("Commit with this message?") {
			fmt.Println("Aborted.")
			return nil
		}
	}

	commitArgs := []string{"commit", "-m", msg}
	if cfg.SignCommits {
		commitArgs = append(commitArgs, "-S")
	}

	gitCmd := exec.Command("git", commitArgs...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}
