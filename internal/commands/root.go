package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	aptfilePath string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "apt-bundle",
	Short: "A declarative package manager for apt",
	Long: `apt-bundle provides a simple, declarative, and shareable way to manage
apt packages and repositories on Debian-based systems, inspired by Homebrew's
brew bundle.`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&aptfilePath, "file", "f", "Aptfile", "Path to Aptfile")
}

// checkRoot verifies that the program is running with root privileges
func checkRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command requires root privileges. Please run with sudo")
	}
	return nil
}
