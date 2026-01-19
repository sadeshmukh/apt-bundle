package apt

import (
	"fmt"
	"os/exec"
)

// lookPath is the function used to look up command paths (overridable for testing)
var lookPath = exec.LookPath

// AddPPA adds a PPA repository using add-apt-repository
func AddPPA(ppa string) error {
	fmt.Printf("Adding PPA: %s\n", ppa)

	if _, err := lookPath("add-apt-repository"); err != nil {
		return fmt.Errorf("add-apt-repository not found. Please install software-properties-common")
	}

	if err := runCommand("add-apt-repository", "-y", ppa); err != nil {
		return wrapCommandError(err, "add PPA", ppa)
	}

	fmt.Printf("✓ PPA %s added successfully\n", ppa)
	return nil
}

// SetLookPath sets the lookPath function (for testing only)
func SetLookPath(f func(string) (string, error)) {
	lookPath = f
}

// ResetLookPath resets lookPath to the default (for testing only)
func ResetLookPath() {
	lookPath = exec.LookPath
}

// AddDebRepository adds a deb repository line to /etc/apt/sources.list.d/
func AddDebRepository(repoLine string) error {
	fmt.Printf("Adding deb repository: %s\n", repoLine)

	// TODO: Implement deb repository addition
	// 1. Generate a filename from hash of the repo line
	// 2. Write the repo line to /etc/apt/sources.list.d/<filename>.list
	// 3. Ensure idempotency

	return nil
}
