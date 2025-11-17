package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.AptfilePath != "Aptfile" {
		t.Errorf("DefaultConfig().AptfilePath = %v, want 'Aptfile'", cfg.AptfilePath)
	}

	if cfg.Verbose != false {
		t.Errorf("DefaultConfig().Verbose = %v, want false", cfg.Verbose)
	}
}

func TestConfigStruct(t *testing.T) {
	// Test custom config creation
	cfg := &Config{
		AptfilePath: "/path/to/custom/Aptfile",
		Verbose:     true,
	}

	if cfg.AptfilePath != "/path/to/custom/Aptfile" {
		t.Errorf("Config.AptfilePath = %v, want '/path/to/custom/Aptfile'", cfg.AptfilePath)
	}

	if cfg.Verbose != true {
		t.Errorf("Config.Verbose = %v, want true", cfg.Verbose)
	}
}

func TestConfigModification(t *testing.T) {
	cfg := DefaultConfig()

	// Test modifying the config
	cfg.AptfilePath = "/custom/path"
	cfg.Verbose = true

	if cfg.AptfilePath != "/custom/path" {
		t.Errorf("After modification, AptfilePath = %v, want '/custom/path'", cfg.AptfilePath)
	}

	if cfg.Verbose != true {
		t.Errorf("After modification, Verbose = %v, want true", cfg.Verbose)
	}
}

func TestConfigIndependence(t *testing.T) {
	// Test that multiple DefaultConfig calls return independent instances
	cfg1 := DefaultConfig()
	cfg2 := DefaultConfig()

	cfg1.Verbose = true
	cfg1.AptfilePath = "/custom1"

	if cfg2.Verbose == true {
		t.Error("Modifying cfg1 should not affect cfg2")
	}

	if cfg2.AptfilePath != "Aptfile" {
		t.Error("Modifying cfg1 should not affect cfg2")
	}
}
