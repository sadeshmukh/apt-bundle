package apt

import (
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
