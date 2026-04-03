package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Summarize staged changes in plain English",
	Long:  `Analyzes staged changes and provides a human-readable summary of what changed and why.`,
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	spinner := ui.NewSpinner("Analyzing changes...")
	spinner.Start()

	summary, err := generator.GenerateDiffSummary(cfg)
	spinner.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", summary)
	return nil
}
