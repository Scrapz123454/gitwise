package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via ldflags at build time.
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of gitwise",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitwise %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
