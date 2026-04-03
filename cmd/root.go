package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitwise",
	Short: "AI-powered commit messages and PR descriptions from your terminal",
	Long: `gitwise analyzes your staged git diffs and generates high-quality
conventional commit messages and pull request descriptions using local or cloud LLMs.

Never write a commit message or PR description by hand again.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("provider", "p", "", "LLM provider (ollama, openai, anthropic, gemini)")
	rootCmd.PersistentFlags().StringP("model", "m", "", "Model name to use")
}
