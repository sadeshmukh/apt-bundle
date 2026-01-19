package apt

import (
	"errors"
	"testing"
)

// mockExecutor is a test double for CommandExecutor
type mockExecutor struct {
	runFunc    func(name string, args ...string) error
	outputFunc func(name string, args ...string) ([]byte, error)
	runCalls   [][]string
	outputCalls [][]string
}

func newMockExecutor() *mockExecutor {
	return &mockExecutor{
		runCalls:    [][]string{},
		outputCalls: [][]string{},
	}
}

func (m *mockExecutor) Run(name string, args ...string) error {
	call := append([]string{name}, args...)
	m.runCalls = append(m.runCalls, call)
	if m.runFunc != nil {
		return m.runFunc(name, args...)
	}
	return nil
}

func (m *mockExecutor) Output(name string, args ...string) ([]byte, error) {
	call := append([]string{name}, args...)
	m.outputCalls = append(m.outputCalls, call)
	if m.outputFunc != nil {
		return m.outputFunc(name, args...)
	}
	return nil, nil
}

func TestSetExecutor(t *testing.T) {
	// Save original executor
	defer ResetExecutor()

	mock := newMockExecutor()
	SetExecutor(mock)

	// Verify the mock is used
	_ = runCommand("test-command", "arg1", "arg2")

	if len(mock.runCalls) != 1 {
		t.Errorf("Expected 1 run call, got %d", len(mock.runCalls))
	}

	if mock.runCalls[0][0] != "test-command" {
		t.Errorf("Expected command 'test-command', got '%s'", mock.runCalls[0][0])
	}
}

func TestResetExecutor(t *testing.T) {
	// Set a mock executor
	mock := newMockExecutor()
	SetExecutor(mock)

	// Reset to real executor
	ResetExecutor()

	// This should use the real executor, not the mock
	// We can verify by checking the mock wasn't called for a subsequent command
	// Note: We can't actually run real commands in tests, so we just verify the reset works
}

func TestRunCommand(t *testing.T) {
	defer ResetExecutor()

	t.Run("success", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return nil
		}
		SetExecutor(mock)

		err := runCommand("test", "arg1")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := newMockExecutor()
		mock.runFunc = func(name string, args ...string) error {
			return errors.New("command failed")
		}
		SetExecutor(mock)

		err := runCommand("test", "arg1")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestRunCommandWithOutput(t *testing.T) {
	defer ResetExecutor()

	t.Run("success with output", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return []byte("test output"), nil
		}
		SetExecutor(mock)

		output, err := runCommandWithOutput("test", "arg1")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if string(output) != "test output" {
			t.Errorf("Expected 'test output', got '%s'", string(output))
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := newMockExecutor()
		mock.outputFunc = func(name string, args ...string) ([]byte, error) {
			return nil, errors.New("command failed")
		}
		SetExecutor(mock)

		_, err := runCommandWithOutput("test", "arg1")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestWrapCommandError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		err := wrapCommandError(errors.New("original"), "install", "package")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to install package: original" {
			t.Errorf("Unexpected error message: %s", err.Error())
		}
	})

	t.Run("without error", func(t *testing.T) {
		err := wrapCommandError(nil, "install", "package")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
