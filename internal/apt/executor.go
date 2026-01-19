package apt

import (
	"fmt"
	"os/exec"
)

type CommandExecutor interface {
	Run(name string, args ...string) error
	Output(name string, args ...string) ([]byte, error)
}

type realExecutor struct{}

func (e *realExecutor) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (e *realExecutor) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

var defaultExecutor CommandExecutor = &realExecutor{}

func runCommand(name string, args ...string) error {
	return defaultExecutor.Run(name, args...)
}

func runCommandWithOutput(name string, args ...string) ([]byte, error) {
	return defaultExecutor.Output(name, args...)
}

func wrapCommandError(err error, operation, target string) error {
	if err != nil {
		return fmt.Errorf("failed to %s %s: %w", operation, target, err)
	}
	return nil
}

// SetExecutor sets the command executor (for testing only)
func SetExecutor(e CommandExecutor) {
	defaultExecutor = e
}

// ResetExecutor resets to the default real executor (for testing only)
func ResetExecutor() {
	defaultExecutor = &realExecutor{}
}
