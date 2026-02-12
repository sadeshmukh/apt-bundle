package commands

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/apt-bundle/apt-bundle/internal/apt"
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

// Mock-based tests for reliable unit testing without root privileges

// mockExecutor implements apt.CommandExecutor for testing
type mockExecutor struct {
	runFunc     func(name string, args ...string) error
	outputFunc  func(name string, args ...string) ([]byte, error)
	runCalls    [][]string
	outputCalls [][]string
}

func newMockExecutor() *mockExecutor {
	return &mockExecutor{
		runCalls:    [][]string{},
		outputCalls: [][]string{},
	}
}

func (m *mockExecutor) Run(name string, args ...string) error {
	call := append([]string{name}, args...)
	m.runCalls = append(m.runCalls, call)
	if m.runFunc != nil {
		return m.runFunc(name, args...)
	}
	return nil
}

func (m *mockExecutor) Output(name string, args ...string) ([]byte, error) {
	call := append([]string{name}, args...)
	m.outputCalls = append(m.outputCalls, call)
	if m.outputFunc != nil {
		return m.outputFunc(name, args...)
	}
	return nil, nil
}

// setupMockRoot sets up the mock to bypass root check
func setupMockRoot() func() {
	SetGetEuid(func() int { return 0 }) // Pretend we're root
	return func() {
		ResetGetEuid()
	}
}

func TestRunInstallWithMock(t *testing.T) {
	t.Run("nonexistent aptfile", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()
		defer apt.ResetExecutor()

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = "/nonexistent/path/Aptfile"

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("Expected error for nonexistent Aptfile")
		}
	})

	t.Run("invalid aptfile", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()
		defer apt.ResetExecutor()

		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "invalid-directive value\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("Expected error for invalid Aptfile")
		}
	})

	t.Run("empty aptfile - no packages", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		// Save and restore noUpdate flag
		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = true // Skip apt-get update

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "# Just a comment\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("successful install with update", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			// Simulate package not installed
			return nil, errors.New("package not found")
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		// Enable update
		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = false

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\napt git\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify apt-get update was called
		updateCalled := false
		for _, call := range mock.runCalls {
			if len(call) >= 2 && call[0] == "apt-get" && call[1] == "update" {
				updateCalled = true
				break
			}
		}
		if !updateCalled {
			t.Error("Expected apt-get update to be called")
		}
	})

	t.Run("install with --no-update", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return nil, errors.New("package not found")
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		// Disable update
		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = true

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify apt-get update was NOT called
		for _, call := range mock.runCalls {
			if len(call) >= 2 && call[0] == "apt-get" && call[1] == "update" {
				t.Error("apt-get update should not be called with --no-update")
			}
		}
	})

	t.Run("package already installed - skip install", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			// Simulate package is installed
			return []byte("install ok installed"), nil
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = true

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify apt-get install was NOT called (package already installed)
		for _, call := range mock.runCalls {
			if len(call) >= 2 && call[0] == "apt-get" && call[1] == "install" {
				t.Error("apt-get install should not be called for already installed package")
			}
		}
	})

	t.Run("update fails", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			if name == "apt-get" && len(args) > 0 && args[0] == "update" {
				return errors.New("E: Could not get lock")
			}
			return nil
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = false

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("Expected error when update fails")
		}
	})

	t.Run("package install fails", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			if name == "apt-get" && len(args) > 0 && args[0] == "install" {
				return errors.New("E: Unable to locate package")
			}
			return nil
		}
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return nil, errors.New("package not found")
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = true

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt nonexistent-package\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err == nil {
			t.Error("Expected error when package install fails")
		}
	})

	t.Run("check installed returns error - warning printed but continues", func(t *testing.T) {
		cleanup := setupMockRoot()
		defer cleanup()

		checkCalls := 0
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			checkCalls++
			// Return error (not installed) but the code should continue
			return nil, errors.New("dpkg-query failed")
		}
		apt.SetExecutor(mock)
		defer apt.ResetExecutor()

		// Setup state path
		tmpDir := t.TempDir()
		apt.SetStatePath(filepath.Join(tmpDir, "state.json"))
		defer apt.ResetStatePath()

		originalNoUpdate := noUpdate
		defer func() { noUpdate = originalNoUpdate }()
		noUpdate = true

		tmpFile := filepath.Join(tmpDir, "Aptfile")
		content := "apt curl\n"

		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()
		aptfilePath = tmpFile

		err := runInstall(installCmd, []string{})
		if err != nil {
			t.Errorf("Expected no error (warning only), got %v", err)
		}

		// Verify IsPackageInstalled was called
		if checkCalls == 0 {
			t.Error("Expected IsPackageInstalled to be called")
		}
	})
}

func TestNoUpdateFlag(t *testing.T) {
	t.Run("flag exists", func(t *testing.T) {
		flag := rootCmd.PersistentFlags().Lookup("no-update")
		if flag == nil {
			t.Fatal("--no-update flag not found")
		}

		if flag.DefValue != "false" {
			t.Errorf("--no-update default = %v, want 'false'", flag.DefValue)
		}
	})
}
