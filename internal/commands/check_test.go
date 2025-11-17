package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckCmd(t *testing.T) {
	t.Run("check command exists", func(t *testing.T) {
		if checkCmd == nil {
			t.Fatal("checkCmd is nil")
		}

		if checkCmd.Use != "check" {
			t.Errorf("checkCmd.Use = %v, want 'check'", checkCmd.Use)
		}

		if checkCmd.RunE == nil {
			t.Error("checkCmd.RunE is nil")
		}
	})
}

func TestRunCheck(t *testing.T) {
	t.Run("with valid aptfile", func(t *testing.T) {
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

		err := runCheck(checkCmd, []string{})
		if err != nil {
			t.Errorf("runCheck() with valid Aptfile returned error: %v", err)
		}
	})

	t.Run("with nonexistent aptfile", func(t *testing.T) {
		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = "/nonexistent/path/Aptfile"

		err := runCheck(checkCmd, []string{})
		if err == nil {
			t.Error("runCheck() with nonexistent Aptfile should return error")
		}
	})

	t.Run("with invalid aptfile", func(t *testing.T) {
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

		err := runCheck(checkCmd, []string{})
		if err == nil {
			t.Error("runCheck() with invalid Aptfile should return error")
		}
	})

	t.Run("with empty aptfile", func(t *testing.T) {
		// Create an empty temporary Aptfile
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "Aptfile")

		if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Save and restore original aptfilePath
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runCheck(checkCmd, []string{})
		if err != nil {
			t.Errorf("runCheck() with empty Aptfile returned error: %v", err)
		}
	})
}
