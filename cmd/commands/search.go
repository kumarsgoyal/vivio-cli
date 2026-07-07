package commands

import (
	fmt "fmt"
	os "os"
	tabwriter "text/tabwriter"

	cobra "github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for channels",
	Long: `Search for channels by name across all countries and categories.

		Results show the global channel number (NO) which can be used with 'vivio play <NO>'.`,
	Example: `  vivio search channels "BBC"
  		vivio search channels "star sports"
  		vivio play 42`,
}

var searchChannelsCmd = &cobra.Command{
	Use:   "channels <query>",
	Short: "Search for channels by name or advanced filters",
	Long: `Search for channels using name or advanced filter syntax.

		Filter Syntax:
  			country:<code>       Filter by country (e.g. country:IN)
  			category:<name>      Filter by category (e.g. category:sports)
  			quality:<res>        Filter by quality (e.g. quality:1080p)
  			language:<code>      Filter by language (e.g. language:eng)

			Use | for OR logic:   country:IN|US|UK

		Examples:
  			vivio search channels "BBC"
  			vivio search channels "country:IN sports"
  			vivio search channels "country:US|UK quality:1080p news"
  			vivio search channels "category:sports quality:720p"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return err
		}

		results := client.AdvancedSearch(args[0])

		if results.Size() == 0 {
			return fmt.Errorf("no channels found matching %q", args[0])
		}

		// Display results in a table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "NO\tCOUNTRY\tNAME\tCATEGORY\tQUALITY\n")

		channels := sortedChannels(results)
		all := sortedChannels(client.GetAll())

		for i := 0; i < results.Size(); i++ {
			ch := channels[i]
			// Find the channel number in the full sorted list
			num := -1
			for j, fullCh := range all {
				if fullCh.ID == ch.ID {
					num = j + 1
					break
				}
			}

			q := ch.Quality
			if q == "" {
				q = "-"
			}

			numStr := "-"
			if num > 0 {
				numStr = fmt.Sprintf("%d", num)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				numStr,
				ch.Country,
				ch.Name,
				ch.Category,
				q,
			)
		}
		w.Flush()

		fmt.Printf("\nFound %d channel(s)\n", results.Size())

		return nil
	},
}

func init() {
	searchCmd.AddCommand(searchChannelsCmd)
}
