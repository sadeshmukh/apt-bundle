package apt

import (
	"fmt"
)

// IsPackageInstalled checks if a package is installed on the system
func IsPackageInstalled(packageName string) (bool, error) {
	output, err := runCommandWithOutput("dpkg-query", "-W", "-f=${Status}", packageName)
	if err != nil {
		return false, nil
	}

	status := string(output)
	return status == "install ok installed", nil
}

// InstallPackage installs a package using apt-get
func InstallPackage(packageName string) error {
	fmt.Printf("Installing package: %s\n", packageName)

	if err := runCommand("apt-get", "install", "-y", packageName); err != nil {
		return wrapCommandError(err, "install package", packageName)
	}

	fmt.Printf("✓ Package %s installed successfully\n", packageName)
	return nil
}

// Update runs apt-get update
func Update() error {
	fmt.Println("Updating package lists...")

	if err := runCommand("apt-get", "update"); err != nil {
		return wrapCommandError(err, "update package lists", "")
	}

	fmt.Println("✓ Package lists updated")
	return nil
}
