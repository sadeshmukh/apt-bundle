package apt

import (
	"errors"
	"os/exec"
	"testing"
)

func TestIsPackageInstalled(t *testing.T) {
	// This test requires dpkg-query to be available
	// We'll test with a package that's likely to be installed (dpkg itself)
	// and one that's likely not installed

	t.Run("check dpkg package", func(t *testing.T) {
		// dpkg should always be installed on Debian/Ubuntu systems
		if _, err := exec.LookPath("dpkg-query"); err != nil {
			t.Skip("dpkg-query not available, skipping test")
		}

		installed, err := IsPackageInstalled("dpkg")
		if err != nil {
			t.Errorf("IsPackageInstalled(dpkg) returned error: %v", err)
		}

		// On systems with dpkg, dpkg itself should be installed
		if !installed {
			t.Log("Note: dpkg not detected as installed, this may be expected on non-Debian systems")
		}
	})

	t.Run("check nonexistent package", func(t *testing.T) {
		if _, err := exec.LookPath("dpkg-query"); err != nil {
			t.Skip("dpkg-query not available, skipping test")
		}

		installed, err := IsPackageInstalled("definitely-not-a-real-package-12345")
		if err != nil {
			t.Errorf("IsPackageInstalled() returned error: %v", err)
		}

		if installed {
			t.Error("IsPackageInstalled() = true for nonexistent package, want false")
		}
	})

	t.Run("empty package name", func(t *testing.T) {
		if _, err := exec.LookPath("dpkg-query"); err != nil {
			t.Skip("dpkg-query not available, skipping test")
		}

		installed, err := IsPackageInstalled("")
		if err != nil {
			// This is acceptable - empty package name might cause an error
			return
		}

		if installed {
			t.Error("IsPackageInstalled('') should return false for empty package name")
		}
	})
}

func TestInstallPackage(t *testing.T) {
	t.Run("install without sudo", func(t *testing.T) {
		// We can't actually test package installation in unit tests
		// as it requires root privileges, but we can verify the function
		// is callable and returns an appropriate error

		err := InstallPackage("test-package-that-does-not-exist")
		// We expect an error because we don't have sudo privileges
		// or the package doesn't exist
		if err == nil {
			t.Log("Warning: InstallPackage succeeded unexpectedly (might have sudo)")
		}
	})

	t.Run("empty package name", func(t *testing.T) {
		err := InstallPackage("")
		// Should fail for empty package name
		if err == nil {
			t.Error("InstallPackage('') should fail for empty package name")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("update without sudo", func(t *testing.T) {
		// Similar to InstallPackage, we can't actually run apt-get update
		// without privileges, but we can verify the function exists
		err := Update()
		if err == nil {
			t.Log("Warning: Update succeeded unexpectedly (might have sudo)")
		}
	})
}

// Benchmark tests
func BenchmarkIsPackageInstalled(b *testing.B) {
	if _, err := exec.LookPath("dpkg-query"); err != nil {
		b.Skip("dpkg-query not available, skipping benchmark")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = IsPackageInstalled("dpkg")
	}
}

// Mock-based tests for reliable unit testing without system dependencies

func TestIsPackageInstalledWithMock(t *testing.T) {
	defer ResetExecutor()

	t.Run("package is installed", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return []byte("install ok installed"), nil
		}
		SetExecutor(mock)

		installed, err := IsPackageInstalled("curl")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !installed {
			t.Error("Expected installed=true, got false")
		}

		// Verify correct command was called
		if len(mock.outputCalls) != 1 {
			t.Errorf("Expected 1 output call, got %d", len(mock.outputCalls))
		}
		if mock.outputCalls[0][0] != "dpkg-query" {
			t.Errorf("Expected dpkg-query, got %s", mock.outputCalls[0][0])
		}
	})

	t.Run("package is not installed - different status", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return []byte("deinstall ok config-files"), nil
		}
		SetExecutor(mock)

		installed, err := IsPackageInstalled("removed-package")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if installed {
			t.Error("Expected installed=false, got true")
		}
	})

	t.Run("package query fails - command error", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return nil, errors.New("dpkg-query: no packages found matching nonexistent")
		}
		SetExecutor(mock)

		installed, err := IsPackageInstalled("nonexistent")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if installed {
			t.Error("Expected installed=false when command fails, got true")
		}
	})

	t.Run("verifies command arguments", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return []byte("install ok installed"), nil
		}
		SetExecutor(mock)

		_, _ = IsPackageInstalled("test-pkg")

		expectedArgs := []string{"dpkg-query", "-W", "-f=${Status}", "test-pkg"}
		if len(mock.outputCalls[0]) != len(expectedArgs) {
			t.Errorf("Expected %d args, got %d", len(expectedArgs), len(mock.outputCalls[0]))
		}
		for i, arg := range expectedArgs {
			if mock.outputCalls[0][i] != arg {
				t.Errorf("Arg %d: expected %s, got %s", i, arg, mock.outputCalls[0][i])
			}
		}
	})
}

func TestInstallPackageWithMock(t *testing.T) {
	defer ResetExecutor()

	t.Run("successful installation", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		SetExecutor(mock)

		err := InstallPackage("curl")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify correct command was called
		if len(mock.runCalls) != 1 {
			t.Errorf("Expected 1 run call, got %d", len(mock.runCalls))
		}
		expectedArgs := []string{"apt-get", "install", "-y", "curl"}
		for i, arg := range expectedArgs {
			if mock.runCalls[0][i] != arg {
				t.Errorf("Arg %d: expected %s, got %s", i, arg, mock.runCalls[0][i])
			}
		}
	})

	t.Run("installation failure", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return errors.New("E: Unable to locate package nonexistent")
		}
		SetExecutor(mock)

		err := InstallPackage("nonexistent")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to install package nonexistent: E: Unable to locate package nonexistent" {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})
}

func TestUpdateWithMock(t *testing.T) {
	defer ResetExecutor()

	t.Run("successful update", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		SetExecutor(mock)

		err := Update()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify correct command was called
		if len(mock.runCalls) != 1 {
			t.Errorf("Expected 1 run call, got %d", len(mock.runCalls))
		}
		expectedArgs := []string{"apt-get", "update"}
		for i, arg := range expectedArgs {
			if mock.runCalls[0][i] != arg {
				t.Errorf("Arg %d: expected %s, got %s", i, arg, mock.runCalls[0][i])
			}
		}
	})

	t.Run("update failure", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return errors.New("E: Could not get lock /var/lib/apt/lists/lock")
		}
		SetExecutor(mock)

		err := Update()
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to update package lists : E: Could not get lock /var/lib/apt/lists/lock" {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})
}
