package main

import (
	"fmt"
	"os"

	"github.com/apt-bundle/apt-bundle/internal/commands"
)

const version = "0.1.0"

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
