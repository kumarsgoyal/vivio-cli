package commands

import (
	fmt "fmt"

	cobra "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vivio version %s\n", versionInfo.version)
		fmt.Printf("commit: %s\n", versionInfo.commit)
		fmt.Printf("built: %s\n", versionInfo.buildDate)
	},
}

func init() {
}
