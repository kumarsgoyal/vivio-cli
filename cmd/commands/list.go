package commands

import (
	fmt "fmt"
	os "os"
	tabwriter "text/tabwriter"

	cobra "github.com/spf13/cobra"
)

var (
	listCountry  string // --country: filter by country code (e.g., IN, US, UK)
	listCategory string // --category: filter by category (e.g., sports, news)
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse and filter channels",
	Long: `Browse and filter channels from 8000+ free live streams worldwide.

		Filters:
  			--country     2-letter country code (e.g. IN, US, GB, AU)
  			--category    channel category (e.g. sports, news, music, movies)

		Examples:
  			vivio list channels
  			vivio list channels --country=IN
  			vivio list channels --category=sports
  			vivio list channels --country=US --category=news`,
}

var listChannelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "List channels with optional filters",
	Example: `  vivio list channels
  		vivio list channels --country=IN
  		vivio list channels --category=sports
  		vivio list channels --country=IN --category=news
  		vivio play 42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return fmt.Errorf("load channels: %w", err)
		}

		if listCountry != "" && !isValidCountryCode(listCountry) {
			fmt.Printf("Invalid country code: %q\n\n", listCountry)
			fmt.Printf("Run: vivio countries list\n\n")
			fmt.Printf("Hint: %s\n", suggestCountries())
			return fmt.Errorf("invalid country code")
		}

		if listCategory != "" && !isValidCategory(listCategory) {
			fmt.Printf("Invalid category: %q\n\n", listCategory)
			fmt.Printf("Run: vivio categories list\n\n")
			fmt.Printf("Hint: %s\n", suggestCategories())
			return fmt.Errorf("invalid category")
		}

		// Get full sorted list for stable numbering
		allSorted := sortedChannels(client.GetAll())

		// Build a map of channel ID → global number
		numberMap := make(map[string]int)
		for i, ch := range allSorted {
			numberMap[ch.ID] = i + 1
		}

		// Filter after numbering
		channels := client.Filter(listCountry, listCategory, "")
		if channels.Size() == 0 {
			fmt.Println("no channels found")
			return nil
		}

		filtered := sortedChannels(channels)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NO\tCOUNTRY\tNAME\tCATEGORY\tQUALITY")
		for _, ch := range filtered {
			num := numberMap[ch.ID]
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", num, ch.Country, ch.Name, ch.Category, ch.Quality)
		}
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n%d channels — play with: vivio play <NO> or vivio play \"<name>\"\n", len(filtered))
		return nil
	},
}

func init() {
	listChannelsCmd.Flags().StringVar(&listCountry, "country", "", "filter by country code (e.g. IN, US, GB)")
	listChannelsCmd.Flags().StringVar(&listCategory, "category", "", "filter by category (e.g. sports, news, music)")
	listCmd.AddCommand(listChannelsCmd)
}
