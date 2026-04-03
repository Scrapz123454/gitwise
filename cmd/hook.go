package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/gitwise/internal/git"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git hook integration",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install gitwise as a prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := git.InstallHook(); err != nil {
			return err
		}
		fmt.Println("Installed gitwise prepare-commit-msg hook.")
		fmt.Println("Now every `git commit` will auto-generate a message.")
		return nil
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the gitwise prepare-commit-msg hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := git.UninstallHook(); err != nil {
			return err
		}
		fmt.Println("Removed gitwise prepare-commit-msg hook.")
		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
	rootCmd.AddCommand(hookCmd)
}
