package commands

import (
	"fmt"

	"github.com/apt-bundle/apt-bundle/internal/apt"
	"github.com/apt-bundle/apt-bundle/internal/aptfile"
	"github.com/spf13/cobra"
)

var (
	noUpdate bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages and repositories from Aptfile",
	Long: `Read the Aptfile and perform the following operations:
1. Add all specified repositories and keys
2. Run apt-get update (unless --no-update is specified)
3. Install all specified packages`,
	RunE: runInstall,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noUpdate, "no-update", false, "Skip updating package lists before installing")
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

	// Load the state to track managed packages
	state, err := apt.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// TODO: Implement the actual installation logic
	// 1. Process all 'key' directives
	// 2. Process all 'ppa' directives
	// 3. Process all 'deb' directives

	// 4. Run apt-get update (unless --no-update is specified)
	if !noUpdate {
		if err := apt.Update(); err != nil {
			return fmt.Errorf("failed to update package lists: %w", err)
		}
	}

	// 5. Process all 'apt' directives to install packages
	packagesToInstall := []string{}
	for _, entry := range entries {
		if entry.Type == "apt" {
			packagesToInstall = append(packagesToInstall, entry.Value)
		}
	}

	if len(packagesToInstall) > 0 {
		fmt.Printf("Installing %d packages...\n", len(packagesToInstall))
		for _, pkg := range packagesToInstall {
			// Check if already installed
			installed, err := apt.IsPackageInstalled(pkg)
			if err != nil {
				fmt.Printf("Warning: Could not check if %s is installed: %v\n", pkg, err)
			}
			if installed {
				fmt.Printf("✓ Package %s is already installed\n", pkg)
				// Track the package in state even if already installed
				// (user may have installed it manually before using apt-bundle)
				state.AddPackage(pkg)
				continue
			}

			// Install the package
			if err := apt.InstallPackage(pkg); err != nil {
				return fmt.Errorf("failed to install package %s: %w", pkg, err)
			}

			// Track successfully installed package in state
			state.AddPackage(pkg)
		}
		fmt.Println("✓ All packages installed successfully")
	} else {
		fmt.Println("No packages to install")
	}

	// Save the updated state
	if err := state.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}
