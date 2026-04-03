package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a CHANGELOG entry from commits",
	Long:  `Generates a structured CHANGELOG entry from conventional commits between tags.`,
	RunE:  runChangelog,
}

func init() {
	changelogCmd.Flags().String("from", "", "Start tag (default: latest tag)")
	changelogCmd.Flags().String("to", "", "End tag (default: HEAD)")
	rootCmd.AddCommand(changelogCmd)
}

func runChangelog(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	fromTag, _ := cmd.Flags().GetString("from")
	toTag, _ := cmd.Flags().GetString("to")

	if fromTag == "" {
		fromTag = git.LatestTag()
	}
	if toTag == "" {
		toTag = "HEAD"
	}

	spinner := ui.NewSpinner("Generating changelog...")
	spinner.Start()

	changelog, err := generator.GenerateChangelog(cfg, fromTag, toTag)
	spinner.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", changelog)
	return nil
}
