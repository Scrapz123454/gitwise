package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint [range]",
	Short: "Lint commit messages for conventional commit compliance",
	Long: `Validates commit messages against conventional commit format.

Examples:
  gitwise lint                    # lint last 10 commits
  gitwise lint HEAD~5..HEAD       # lint specific range`,
	RunE: runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	commitRange := ""
	if len(args) > 0 {
		commitRange = args[0]
	}

	spinner := ui.NewSpinner("Linting commits...")
	spinner.Start()

	result, err := generator.LintCommits(cfg, commitRange)
	spinner.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", result)
	return nil
}
