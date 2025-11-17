package apt

import (
	"fmt"
	"os/exec"
)

// IsPackageInstalled checks if a package is installed on the system
func IsPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Status}", packageName)
	output, err := cmd.Output()
	if err != nil {
		// Package not found
		return false, nil
	}

	status := string(output)
	return status == "install ok installed", nil
}

// InstallPackage installs a package using apt-get
func InstallPackage(packageName string) error {
	fmt.Printf("Installing package: %s\n", packageName)

	cmd := exec.Command("apt-get", "install", "-y", packageName)
	cmd.Stdout = nil // TODO: Wire up proper output handling
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", packageName, err)
	}

	fmt.Printf("✓ Package %s installed successfully\n", packageName)
	return nil
}

// Update runs apt-get update
func Update() error {
	fmt.Println("Updating package lists...")

	cmd := exec.Command("apt-get", "update")
	cmd.Stdout = nil // TODO: Wire up proper output handling
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update package lists: %w", err)
	}

	fmt.Println("✓ Package lists updated")
	return nil
}
