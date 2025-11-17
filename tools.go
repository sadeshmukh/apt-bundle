//go:build tools
// +build tools

// Package tools tracks tool dependencies for this project
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

