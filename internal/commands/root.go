package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	aptfilePath string
	version     = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "apt-bundle",
	Short: "A declarative package manager for apt",
	Long: `apt-bundle provides a simple, declarative, and shareable way to manage
apt packages and repositories on Debian-based systems, inspired by Homebrew's
brew bundle.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&aptfilePath, "file", "f", "Aptfile", "Path to Aptfile")
}

func checkRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command requires root privileges. Please run with sudo")
	}
	return nil
}
