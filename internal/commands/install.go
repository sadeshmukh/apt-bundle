package commands

import (
	"fmt"

	"github.com/apt-bundle/apt-bundle/internal/aptfile"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages and repositories from Aptfile",
	Long: `Read the Aptfile and perform the following operations:
1. Add all specified repositories and keys
2. Run apt-get update
3. Install all specified packages`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	// Make install the default command
	rootCmd.RunE = runInstall
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Check for root privileges
	if err := checkRoot(); err != nil {
		return err
	}

	fmt.Printf("Reading Aptfile from: %s\n", aptfilePath)

	// Parse the Aptfile
	entries, err := aptfile.Parse(aptfilePath)
	if err != nil {
		return fmt.Errorf("failed to parse Aptfile: %w", err)
	}

	fmt.Printf("Found %d entries in Aptfile\n", len(entries))

	// TODO: Implement the actual installation logic
	// 1. Process all 'key' directives
	// 2. Process all 'ppa' directives
	// 3. Process all 'deb' directives
	// 4. Run apt-get update
	// 5. Process all 'apt' directives to install packages

	return nil
}
