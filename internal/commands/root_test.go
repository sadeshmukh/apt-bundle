package commands

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	t.Run("help flag", func(t *testing.T) {
		// Test that the Execute function can be called
		// We'll test with --help to avoid requiring actual execution
		os.Args = []string{"apt-bundle", "--help"}
		err := Execute()
		// --help causes cobra to print help and return nil or a specific error
		// Either is acceptable
		if err != nil {
			t.Logf("Execute() with --help returned: %v", err)
		}
	})

	t.Run("version flag", func(t *testing.T) {
		os.Args = []string{"apt-bundle", "--version"}
		err := Execute()
		if err != nil {
			t.Logf("Execute() with --version returned: %v", err)
		}
	})
}

func TestRootCmd(t *testing.T) {
	t.Run("root command exists", func(t *testing.T) {
		if rootCmd == nil {
			t.Fatal("rootCmd is nil")
		}

		if rootCmd.Use != "apt-bundle" {
			t.Errorf("rootCmd.Use = %v, want 'apt-bundle'", rootCmd.Use)
		}

		if rootCmd.Version == "" {
			t.Error("rootCmd.Version should not be empty")
		}
	})

	t.Run("global flags", func(t *testing.T) {
		flag := rootCmd.PersistentFlags().Lookup("file")
		if flag == nil {
			t.Fatal("--file flag not found")
		}

		if flag.Shorthand != "f" {
			t.Errorf("--file shorthand = %v, want 'f'", flag.Shorthand)
		}

		if flag.DefValue != "Aptfile" {
			t.Errorf("--file default = %v, want 'Aptfile'", flag.DefValue)
		}
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := rootCmd.Commands()
		if len(commands) == 0 {
			t.Error("rootCmd has no subcommands")
		}

		// Check for expected subcommands
		cmdNames := make(map[string]bool)
		for _, cmd := range commands {
			cmdNames[cmd.Name()] = true
		}

		expectedCommands := []string{"check", "dump", "install"}
		for _, expected := range expectedCommands {
			if !cmdNames[expected] {
				t.Errorf("Expected subcommand %q not found", expected)
			}
		}
	})
}

func TestCheckRoot(t *testing.T) {
	t.Run("check root privileges", func(t *testing.T) {
		err := checkRoot()

		// Get current effective UID
		euid := os.Geteuid()

		if euid == 0 {
			// Running as root - should not error
			if err != nil {
				t.Errorf("checkRoot() returned error when running as root: %v", err)
			}
		} else {
			// Not running as root - should error
			if err == nil {
				t.Error("checkRoot() should return error when not running as root")
			}
		}
	})
}

func TestGlobalVariables(t *testing.T) {
	t.Run("aptfilePath variable", func(t *testing.T) {
		// Save original value
		originalPath := aptfilePath
		defer func() { aptfilePath = originalPath }()

		// Test that we can modify the variable
		testPath := "/custom/path/Aptfile"
		aptfilePath = testPath

		if aptfilePath != testPath {
			t.Errorf("aptfilePath = %v, want %v", aptfilePath, testPath)
		}
	})
}
