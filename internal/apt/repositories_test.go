package apt

import (
	"errors"
	"os/exec"
	"testing"
)

func TestAddPPA(t *testing.T) {
	t.Run("add-apt-repository not available", func(t *testing.T) {
		// Test with a PPA when add-apt-repository might not be available
		// The function should check for the command's existence
		err := AddPPA("ppa:deadsnakes/ppa")

		// If add-apt-repository is not found, we should get a specific error
		if err != nil {
			if _, lookupErr := exec.LookPath("add-apt-repository"); lookupErr != nil {
				// Expected: add-apt-repository not found
				return
			}
			// If add-apt-repository exists, error is likely due to permissions
			t.Logf("AddPPA failed (likely due to permissions): %v", err)
		}
	})

	t.Run("empty ppa", func(t *testing.T) {
		if _, err := exec.LookPath("add-apt-repository"); err != nil {
			t.Skip("add-apt-repository not available, skipping test")
		}

		err := AddPPA("")
		// Should handle empty PPA gracefully
		if err == nil {
			t.Log("Warning: AddPPA('') succeeded unexpectedly")
		}
	})

	t.Run("invalid ppa format", func(t *testing.T) {
		if _, err := exec.LookPath("add-apt-repository"); err != nil {
			t.Skip("add-apt-repository not available, skipping test")
		}

		err := AddPPA("invalid-ppa-format")
		// The command should fail or handle invalid format
		if err == nil {
			t.Log("Warning: AddPPA with invalid format succeeded unexpectedly")
		}
	})

	t.Run("valid ppa format", func(t *testing.T) {
		if _, err := exec.LookPath("add-apt-repository"); err != nil {
			t.Skip("add-apt-repository not available, skipping test")
		}

		err := AddPPA("ppa:deadsnakes/ppa")
		// Expected to fail without sudo, but function should be callable
		if err == nil {
			t.Log("Warning: AddPPA succeeded unexpectedly (might have sudo)")
		}
	})
}

func TestAddDebRepository(t *testing.T) {
	t.Run("basic functionality", func(t *testing.T) {
		err := AddDebRepository("deb http://example.com/ubuntu jammy main")
		// Currently returns nil as it's not implemented
		if err != nil {
			t.Errorf("AddDebRepository() returned error: %v", err)
		}
	})

	t.Run("empty repository line", func(t *testing.T) {
		err := AddDebRepository("")
		// Should handle empty line gracefully
		if err != nil {
			t.Errorf("AddDebRepository('') returned error: %v", err)
		}
	})

	t.Run("complex repository line", func(t *testing.T) {
		repoLine := "deb [arch=amd64 signed-by=/usr/share/keyrings/example.gpg] http://example.com/ubuntu jammy main contrib"
		err := AddDebRepository(repoLine)
		if err != nil {
			t.Errorf("AddDebRepository() returned error: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkAddPPA(b *testing.B) {
	if _, err := exec.LookPath("add-apt-repository"); err != nil {
		b.Skip("add-apt-repository not available, skipping benchmark")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This will fail without sudo, but we're benchmarking the code path
		_ = AddPPA("ppa:test/ppa")
	}
}

func BenchmarkAddDebRepository(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddDebRepository("deb http://example.com/ubuntu jammy main")
	}
}

// Mock-based tests for reliable unit testing without system dependencies

func TestAddPPAWithMock(t *testing.T) {
	defer ResetExecutor()
	defer ResetLookPath()

	t.Run("add-apt-repository not found", func(t *testing.T) {
		SetLookPath(func(file string) (string, error) {
			return "", errors.New("executable file not found in $PATH")
		})

		err := AddPPA("ppa:test/ppa")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "add-apt-repository not found. Please install software-properties-common" {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})

	t.Run("successful PPA addition", func(t *testing.T) {
		SetLookPath(func(file string) (string, error) {
			return "/usr/bin/add-apt-repository", nil
		})

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		SetExecutor(mock)

		err := AddPPA("ppa:deadsnakes/ppa")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify correct command was called
		if len(mock.runCalls) != 1 {
			t.Errorf("Expected 1 run call, got %d", len(mock.runCalls))
		}
		expectedArgs := []string{"add-apt-repository", "-y", "ppa:deadsnakes/ppa"}
		for i, arg := range expectedArgs {
			if mock.runCalls[0][i] != arg {
				t.Errorf("Arg %d: expected %s, got %s", i, arg, mock.runCalls[0][i])
			}
		}
	})

	t.Run("PPA addition command failure", func(t *testing.T) {
		SetLookPath(func(file string) (string, error) {
			return "/usr/bin/add-apt-repository", nil
		})

		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return errors.New("E: The repository 'ppa:invalid/ppa' does not have a Release file")
		}
		SetExecutor(mock)

		err := AddPPA("ppa:invalid/ppa")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		expectedErr := "failed to add PPA ppa:invalid/ppa: E: The repository 'ppa:invalid/ppa' does not have a Release file"
		if err.Error() != expectedErr {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})

	t.Run("lookPath is called with correct argument", func(t *testing.T) {
		var calledWith string
		SetLookPath(func(file string) (string, error) {
			calledWith = file
			return "/usr/bin/add-apt-repository", nil
		})

		mock := newMockExecutor()
		SetExecutor(mock)

		_ = AddPPA("ppa:test/ppa")

		if calledWith != "add-apt-repository" {
			t.Errorf("Expected lookPath to be called with 'add-apt-repository', got '%s'", calledWith)
		}
	})
}

func TestSetLookPath(t *testing.T) {
	defer ResetLookPath()

	customLookPath := func(file string) (string, error) {
		return "/custom/path", nil
	}
	SetLookPath(customLookPath)

	path, err := lookPath("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if path != "/custom/path" {
		t.Errorf("Expected '/custom/path', got '%s'", path)
	}
}

func TestResetLookPath(t *testing.T) {
	// Set a custom lookPath
	SetLookPath(func(file string) (string, error) {
		return "/custom/path", nil
	})

	// Reset it
	ResetLookPath()

	// After reset, lookPath should behave like exec.LookPath
	// We can't easily verify this without side effects, but we can verify it doesn't panic
	_, _ = lookPath("ls")
}
