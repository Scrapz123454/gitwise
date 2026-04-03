package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "AI code review of staged changes",
	Long:  `Performs an AI-powered code review analyzing staged changes for bugs, security issues, and improvements.`,
	RunE:  runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}

func runReview(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	spinner := ui.NewSpinner("Reviewing code...")
	spinner.Start()

	review, err := generator.GenerateReview(cfg)
	spinner.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", review)
	return nil
}
