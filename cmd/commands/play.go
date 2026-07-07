package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"

	core "github.com/viviotv/vivio/pkg"
)

var (
	playerFlag  string // --player: media player to use (mpv or ffplay)
	qualityFlag string // --quality: preferred quality (1080p, 720p, etc.)
)

var playCmd = &cobra.Command{
	Use:   "play <no or channel name>",
	Short: "Play a channel in mpv or ffplay",
	Long: `Play a live TV channel using mpv or ffplay.

		Flags:
  			--player      media player to use: mpv or ffplay (auto-detects if not set)
  			--quality     preferred quality: 1080p, 720p, 576p (uses best available if not set)

		Examples:
  			vivio play 42
  			vivio play "BBC News"
  			vivio play "BBC News" --player=mpv
  			vivio play "BBC News" --quality=1080p
  			vivio play "star sports" --quality=720p

			Note: Channel numbers must be positive. Negative numbers like -5 will show
			a flag error - use positive channel numbers only (1 and above).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initClient(); err != nil {
			return err
		}
		if err := client.Load(); err != nil {
			return fmt.Errorf("load channels: %w", err)
		}

		var ch *core.Channel

		// numeric arg, pick by position from the sorted full list
		if n, err := strconv.Atoi(args[0]); err == nil {
			if n <= 0 {
				return fmt.Errorf("channel number must be positive (valid: 1 and above)")
			}

			all := sortedChannels(client.GetAll())
			if n > len(all) {
				return fmt.Errorf("channel #%d out of range (valid: 1–%d)", n, len(all))
			}
			ch = all[n-1]
		} else {
			// string arg, search by name
			results := client.Search(args[0])
			if results.Size() == 0 {
				return fmt.Errorf("no channel found matching %q", args[0])
			}
			ch = results.Get(0)
		}

		// If user specified quality, filter streams
		if qualityFlag != "" {
			ch = filterByQuality(ch, qualityFlag)
		}

		label := fmt.Sprintf("%s [%s] %s", ch.Name, ch.Country, ch.Quality)
		fmt.Printf("Playing: %s\n", label)

		// Try primary stream first, fallback to alternatives on error
		return playWithFallback(ch, label, playerFlag)
	},
}

func init() {
	playCmd.Flags().StringVar(&playerFlag, "player", "", "media player to use: mpv or ffplay (auto-detects if not set)")
	playCmd.Flags().StringVar(&qualityFlag, "quality", "", "preferred quality: 1080p, 720p, 576p (uses best available if not set)")
}

// filterByQuality filters channel streams to match requested quality
// Reorders streams to prioritize the requested quality
func filterByQuality(ch *core.Channel, quality string) *core.Channel {
	if len(ch.Streams) == 0 {
		return ch
	}

	// Find streams matching requested quality
	var matchingStreams []core.Stream
	var otherStreams []core.Stream

	for _, stream := range ch.Streams {
		if stream.Quality == quality {
			matchingStreams = append(matchingStreams, stream)
		} else {
			otherStreams = append(otherStreams, stream)
		}
	}

	// If no exact match, return original
	if len(matchingStreams) == 0 {
		available := getAvailableQualities(ch)
		fmt.Printf("Note: %s quality not available.\n", quality)
		if len(available) > 0 {
			fmt.Printf("      Available qualities: %s\n", joinQualities(available))
			fmt.Printf("      Using: %s\n", ch.Quality)
		} else {
			fmt.Printf("      Using best available.\n")
		}
		return ch
	}

	// Put matching quality streams first, then others as fallback
	filteredCh := *ch
	filteredCh.Streams = make([]core.Stream, 0, len(matchingStreams)+len(otherStreams))
	filteredCh.Streams = append(filteredCh.Streams, matchingStreams...)
	filteredCh.Streams = append(filteredCh.Streams, otherStreams...)
	filteredCh.StreamURL = filteredCh.Streams[0].URL
	filteredCh.Quality = filteredCh.Streams[0].Quality

	fmt.Printf("Note: Found %d stream(s) with %s quality.\n", len(matchingStreams), quality)
	return &filteredCh
}

// playWithFallback tries primary stream first, then fallbacks on error
func playWithFallback(ch *core.Channel, label, player string) error {
	if len(ch.Streams) == 0 {
		return fmt.Errorf("no streams available for %s", ch.Name)
	}

	// If only one stream, just play it (no fallback)
	if len(ch.Streams) == 1 {
		return launch(ch.Streams[0].URL, player)
	}

	// Multiple streams available - try with fallback
	fmt.Printf("Note: %d streams available. Will try alternatives if primary fails.\n\n", len(ch.Streams))

	for i, stream := range ch.Streams {
		streamNum := ""
		if i > 0 {
			streamNum = fmt.Sprintf(" (fallback %d/%d)", i, len(ch.Streams)-1)
		}

		qualityLabel := ""
		if stream.Quality != "" {
			qualityLabel = fmt.Sprintf(" [%s]", stream.Quality)
		}

		warningLabel := ""
		if stream.Label != "" {
			warningLabel = fmt.Sprintf(" ⚠ %s", stream.Label)
		}

		if i > 0 {
			fmt.Printf("\n━━━ Trying alternative stream%s%s%s...\n", streamNum, qualityLabel, warningLabel)
		}

		err := launch(stream.URL, player)
		if err == nil {
			return nil // Success - player exited cleanly
		}

		// Stream failed - if not the last one, try next
		if i < len(ch.Streams)-1 {
			fmt.Printf("━━━ Stream failed. Trying next alternative...\n")
			continue
		}

		// Last stream also failed
		return fmt.Errorf("all %d streams failed for %s. Last error: %w", len(ch.Streams), ch.Name, err)
	}

	return nil
}

// launch opens the stream URL in the specified player.
// If player is empty, tries mpv then ffplay in order.
func launch(url, player string) error {
	type playerDef struct {
		bin  string
		args []string
	}

	all := map[string]playerDef{
		"mpv":    {"mpv", []string{"--title=Vivio", url}},
		"ffplay": {"ffplay", []string{"-i", url, "-window_title", "Vivio"}},
	}

	if player != "" {
		p, ok := all[player]
		if !ok {
			return fmt.Errorf("unknown player %q — use mpv or ffplay", player)
		}
		path, err := exec.LookPath(p.bin)
		if err != nil {
			return fmt.Errorf("%s is not installed — run: sudo dnf install %s", player, player)
		}
		return run(path, p.args...)
	}

	for _, name := range []string{"mpv", "ffplay"} {
		p := all[name]
		path, err := exec.LookPath(p.bin)
		if err != nil {
			continue
		}
		fmt.Printf("Using: %s\n", name)
		return run(path, p.args...)
	}

	fmt.Println("No player found. Install mpv:  sudo dnf install mpv")
	fmt.Printf("Or copy this URL into VLC:\n%s\n", url)
	return nil
}

func run(path string, args ...string) error {
	c := exec.Command(path, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
