package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	core "github.com/viviotv/vivio/pkg"
)

var infoCmd = &cobra.Command{
	Use:   "info <channel number or name>",
	Short: "Show detailed information about a channel",
	Long: `Display complete channel details including:
  		- Name, ID, Country, Category, Language
  		- Network, Website, Logo
  		- All available streams with quality indicators
  		- Stream URLs, geo-blocking info, user-agent requirements

		Examples:
  			vivio info 3736
  			vivio info "BBC News"
  			vivio info "Star Sports"

		Note: Channel numbers must be positive. Negative numbers like -5 will show
		a flag error - use positive channel numbers only (1 and above).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return err
		}

		ch, err := findChannel(args[0])
		if err != nil {
			return err
		}

		printChannelInfo(ch)
		return nil
	},
}

func init() {
}

// findChannel finds a channel by number or name
func findChannel(query string) (*core.Channel, error) {
	// Try as number first
	if n, err := strconv.Atoi(query); err == nil {
		if n <= 0 {
			return nil, fmt.Errorf("channel number must be positive (valid: 1 and above)")
		}

		all := sortedChannels(client.GetAll())
		if n > len(all) {
			return nil, fmt.Errorf("channel #%d out of range (valid: 1–%d)", n, len(all))
		}
		return all[n-1], nil
	}

	// Search by name
	results := client.Search(query)
	if results.Size() == 0 {
		return nil, fmt.Errorf("no channel found matching %q", query)
	}
	if results.Size() > 1 {
		fmt.Printf("Multiple channels found. Showing first match:\n\n")
	}
	return results.Get(0), nil
}

// printChannelInfo prints all channel information
func printChannelInfo(ch *core.Channel) {
	printChannelHeader(ch)
	printChannelDetails(ch)
	printAvailableQualities(ch)
	printStreams(ch)
	printQuickActions(ch)
}

// printChannelHeader prints the channel name header
func printChannelHeader(ch *core.Channel) {
	fmt.Printf("\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  %s\n", ch.Name)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
}

// printChannelDetails prints channel metadata
func printChannelDetails(ch *core.Channel) {
	fmt.Printf("  📺 Channel Details\n")
	fmt.Printf("  ├─ ID:       %s\n", ch.ID)
	fmt.Printf("  ├─ Country:  %s (%s)\n", core.CountryName(ch.Country), ch.Country)
	fmt.Printf("  ├─ Category: %s\n", ch.Category)
	if ch.Language != "" {
		fmt.Printf("  ├─ Language: %s\n", ch.Language)
	}
	if ch.Network != "" {
		fmt.Printf("  ├─ Network:  %s\n", ch.Network)
	}
	if ch.Website != "" {
		fmt.Printf("  ├─ Website:  %s\n", ch.Website)
	}
	fmt.Printf("  └─ Logo:     %s\n", ch.Logo)
}

// printAvailableQualities prints quality summary
func printAvailableQualities(ch *core.Channel) {
	qualities := getAvailableQualities(ch)
	if len(qualities) > 0 {
		fmt.Printf("\n  🎬 Available Qualities: %s\n", joinQualities(qualities))
		fmt.Printf("     Use: vivio play %q --quality=<quality>\n", ch.Name)
	}
}

// printStreams prints all available streams grouped by quality
func printStreams(ch *core.Channel) {
	qualities := getAvailableQualities(ch)
	streamsByQuality := groupStreamsByQuality(ch)

	fmt.Printf("\n  📡 Available Streams (%d total)\n", len(ch.Streams))

	for _, quality := range qualities {
		streams := streamsByQuality[quality]
		if len(streams) == 0 {
			continue
		}

		fmt.Printf("\n  ┌─ %s (%d stream%s)\n", quality, len(streams), pluralize(len(streams)))

		for i := range streams {
			printStream(&streams[i], &ch.Streams[0], i, len(streams))
		}
	}
}

func groupStreamsByQuality(ch *core.Channel) map[string][]core.Stream {
	groups := make(map[string][]core.Stream)

	for _, stream := range ch.Streams {
		quality := stream.Quality
		if quality == "" {
			quality = "Unknown"
		}
		groups[quality] = append(groups[quality], stream)
	}

	return groups
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// printStream prints a single stream with formatting
func printStream(stream, primaryStream *core.Stream, index, total int) {
	labelStr := ""
	if stream.Label != "" {
		labelStr = fmt.Sprintf(" ⚠ %s", stream.Label)
	}

	primary := ""
	if stream.URL == primaryStream.URL {
		primary = " ★"
	}

	prefix := "├"
	if index == total-1 {
		prefix = "└"
	}

	fmt.Printf("  %s─ %s%s%s\n", prefix, stream.URL, labelStr, primary)

	if stream.UserAgent != "" {
		fmt.Printf("     │  User-Agent: %s\n", stream.UserAgent)
	}
	if stream.Referrer != "" {
		fmt.Printf("     │  Referrer: %s\n", stream.Referrer)
	}
}

// printQuickActions prints quick action commands
func printQuickActions(ch *core.Channel) {
	qualities := getAvailableQualities(ch)

	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  💡 Quick Actions\n")
	fmt.Printf("  ─────────────────────────────────────────────────\n")
	fmt.Printf("  Play:         vivio play %q\n", ch.Name)
	if len(qualities) > 0 {
		fmt.Printf("  With quality: vivio play %q --quality=%s\n", ch.Name, qualities[0])
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
}
