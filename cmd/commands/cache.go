package commands

import (
	fmt "fmt"

	cobra "github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage local channel cache",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete the local channel cache (next command will re-fetch)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.ClearCache(); err != nil {
			return err
		}
		fmt.Println("cache cleared")
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
}
