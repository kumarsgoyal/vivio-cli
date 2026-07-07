package commands

import (
	fmt "fmt"
	sort "sort"

	cobra "github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Browse channels by category",
}

var categoriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available categories",
	Example: `  vivio categories list
		vivio list --category=sports
		vivio list --category=news
		vivio list --country=IN --category=sports`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return fmt.Errorf("load channels: %w", err)
		}
		list := client.Categories()

		items := make([]string, 0, list.Size())
		for i := 0; i < list.Size(); i++ {
			items = append(items, list.Get(i))
		}
		sort.Strings(items)

		for _, c := range items {
			fmt.Println(c)
		}
		return nil
	},
}

func init() {
	categoriesCmd.AddCommand(categoriesListCmd)
}
