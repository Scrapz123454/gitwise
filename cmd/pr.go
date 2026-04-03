package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aymenhmaidiwastaken/gitwise/internal/generator"
	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/aymenhmaidiwastaken/gitwise/internal/ui"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a PR description for the current branch",
	Long:  `Analyzes all commits on the current branch and generates a structured PR description.`,
	RunE:  runPR,
}

func init() {
	prCmd.Flags().StringP("base", "b", "", "Base branch to compare against (default: auto-detect)")
	prCmd.Flags().Bool("create", false, "Create the PR on GitHub using gh CLI")
	prCmd.Flags().Bool("labels", false, "Suggest labels for the PR")
	rootCmd.AddCommand(prCmd)
}

func runPR(cmd *cobra.Command, args []string) error {
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not inside a git repository")
	}

	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	base, _ := cmd.Flags().GetString("base")
	create, _ := cmd.Flags().GetBool("create")
	showLabels, _ := cmd.Flags().GetBool("labels")

	spinner := ui.NewSpinner("Generating PR description...")
	spinner.Start()

	result, err := generator.GeneratePRDescription(cfg, base)
	spinner.Stop()
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", result.Description)
	if cfg.ShowCost && result.CostEstimate != "" {
		fmt.Printf("\n%s\n", ui.FormatCostInfo(result.TokenCount, result.CostEstimate))
	}

	// Suggest labels
	if showLabels {
		fmt.Println()
		labelSpinner := ui.NewSpinner("Suggesting labels...")
		labelSpinner.Start()
		labels, err := generator.SuggestLabels(cfg, base)
		labelSpinner.Stop()
		if err == nil && len(labels) > 0 {
			fmt.Printf("Suggested labels: %s\n", strings.Join(labels, ", "))
			result.Labels = labels
		}
	}

	if create {
		return createGitHubPR(result, base)
	}

	return nil
}

func createGitHubPR(result *generator.PRResult, base string) error {

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed — install it from https://cli.github.com")
	}

	branch, _ := git.CurrentBranch()
	if base == "" {
		base = git.DefaultBranch()
	}

	title := branch
	lines := strings.SplitN(result.Description, "\n", 5)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			title = line
			break
		}
	}

	if !ui.Confirm(fmt.Sprintf("Create PR '%s' -> %s?", branch, base)) {
		fmt.Println("Aborted.")
		return nil
	}

	ghArgs := []string{"pr", "create",
		"--title", title,
		"--body", result.Description,
		"--base", base,
	}

	// Add labels if available
	if len(result.Labels) > 0 {
		ghArgs = append(ghArgs, "--label", strings.Join(result.Labels, ","))
	}

	ghCmd := exec.Command("gh", ghArgs...)
	ghCmd.Stdout = os.Stdout
	ghCmd.Stderr = os.Stderr

	if err := ghCmd.Run(); err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	return nil
}
