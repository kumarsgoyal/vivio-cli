package commands

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	core "github.com/viviotv/vivio/pkg"
)

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "Browse channels by country",
}

var countriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all countries that have active channels",
	Example: `  vivio countries list
  		vivio list --country=IN
  		vivio list --country=IN --category=sports
  		vivio play "Star Sports" --country=IN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return fmt.Errorf("load channels: %w", err)
		}

		list := client.Countries()

		type row struct {
			code string
			name string
		}

		rows := make([]row, 0, list.Size())
		for i := 0; i < list.Size(); i++ {
			code := list.Get(i)
			rows = append(rows, row{code: code, name: core.CountryName(code)})
		}

		sort.Slice(rows, func(i, j int) bool {
			return rows[i].name < rows[j].name
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "CODE\tCOUNTRY")
		for _, r := range rows {
			fmt.Fprintf(w, "%s\t%s\n", r.code, r.name)
		}
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n%d countries with active streams\n", len(rows))
		return nil
	},
}

func init() {
	countriesCmd.AddCommand(countriesListCmd)
}
