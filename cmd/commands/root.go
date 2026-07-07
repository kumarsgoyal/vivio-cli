package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	core "github.com/viviotv/vivio/pkg"
)

var client *core.CoreClient

var rootCmd = &cobra.Command{
	Use:   "vivio",
	Short: "Free live TV — 8000+ channels worldwide, zero cost",
	Long: `Vivio streams 8000+ free public TV channels from around the world.
		Channels come directly from iptv-org. No account needed. No ads.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	versionInfo = struct {
		version   string
		commit    string
		buildDate string
	}{
		version:   "dev",
		commit:    "unknown",
		buildDate: "unknown",
	}
)

// SetVersion sets the version information
func SetVersion(version, commit, buildDate string) {
	versionInfo.version = version
	versionInfo.commit = commit
	versionInfo.buildDate = buildDate
}

func initClient() error {
	cacheDir, err := defaultCacheDir()
	if err != nil {
		return err
	}
	client = core.NewCoreClient(cacheDir)
	return nil
}

func defaultCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	dir := filepath.Join(home, ".cache", "vivio")
	return dir, os.MkdirAll(dir, 0o755)
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(countriesCmd)
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(versionCmd)
}
