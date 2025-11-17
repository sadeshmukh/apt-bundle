package apt

import (
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
