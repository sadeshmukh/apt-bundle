package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCmd(t *testing.T) {
	t.Run("install command exists", func(t *testing.T) {
		if installCmd == nil {
			t.Fatal("installCmd is nil")
		}

		if installCmd.Use != "install" {
			t.Errorf("installCmd.Use = %v, want 'install'", installCmd.Use)
		}

		if installCmd.RunE == nil {
			t.Error("installCmd.RunE is nil")
		}
	})

	t.Run("install is default command", func(t *testing.T) {
		// Check that rootCmd has a RunE function (install as default)
		if rootCmd.RunE == nil {
			t.Error("rootCmd.RunE should be set to make install the default command")
		}
	})
}

func TestRunInstall(t *testing.T) {
	t.Run("without root privileges", func(t *testing.T) {
		// Skip if running as root
		if os.Geteuid() == 0 {
			t.Skip("Skipping test - running as root")
		}

		// Create a temporary Aptfile
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("runInstall() should fail without root privileges")
		}
	})

	t.Run("with nonexistent aptfile as root", func(t *testing.T) {
		// Only run if we're root
		if os.Geteuid() != 0 {
			t.Skip("Skipping test - requires root privileges")
		}

		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = "/nonexistent/path/Aptfile"

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("runInstall() with nonexistent Aptfile should return error")
		}
	})

	t.Run("with invalid aptfile as root", func(t *testing.T) {
		// Only run if we're root
		if os.Geteuid() != 0 {
			t.Skip("Skipping test - requires root privileges")
		}

		// Create a temporary Aptfile with invalid content
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "invalid-directive value\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("runInstall() with invalid Aptfile should return error")
		}
	})

	t.Run("with valid aptfile as root", func(t *testing.T) {
		// Only run if we're root
		if os.Geteuid() != 0 {
			t.Skip("Skipping test - requires root privileges")
		}

		// Create a temporary Aptfile
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\napt git\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("runInstall() with valid Aptfile returned error: %v", err)
		}
	})
}
