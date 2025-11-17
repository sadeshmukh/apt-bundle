package config

// Config holds the application configuration
type Config struct {
	AptfilePath string
	Verbose     bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		AptfilePath: "Aptfile",
		Verbose:     false,
	}
}
