# Testing Guide

## Quick Start

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Generate HTML Coverage Report
```bash
make test-coverage-html
# Open coverage.html in your browser
```

## Test Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests with verbose output |
| `make test-coverage` | Run tests with coverage report in terminal |
| `make test-coverage-html` | Generate HTML coverage report |
| `make fmt` | Format all Go code |
| `make vet` | Run go vet static analysis |
| `make lint` | Run golangci-lint |

## Coverage Summary

- **Overall**: 85.7%
- **internal/config**: 100.0%
- **internal/aptfile**: 98.3%
- **internal/apt**: 80.6%
- **internal/commands**: 75.0%

## Test Files

```
internal/
├── apt/
│   ├── keys_test.go          (GPG key management tests)
│   ├── packages_test.go      (Package management tests)
│   └── repositories_test.go  (Repository management tests)
├── aptfile/
│   └── parser_test.go        (Aptfile parsing tests)
├── commands/
│   ├── check_test.go         (Check command tests)
│   ├── dump_test.go          (Dump command tests)
│   ├── install_test.go       (Install command tests)
│   └── root_test.go          (Root command tests)
└── config/
    └── config_test.go        (Configuration tests)
```

## Test Features

✅ **Comprehensive Coverage**: All major functions are tested  
✅ **Edge Cases**: Empty inputs, invalid inputs, error conditions  
✅ **Race Detection**: Tests run with `-race` flag  
✅ **Benchmarks**: Performance benchmarks included  
✅ **Table-Driven**: Clean, maintainable test structure  
✅ **Isolated**: Tests use temporary directories  

## CI/CD Ready

The test suite is optimized for CI/CD:
- Fast execution (< 5 seconds)
- Handles permissions gracefully
- Produces standard coverage reports
- Proper exit codes

## Writing New Tests

### Test File Template

```go
package yourpackage

import (
    "testing"
)

func TestYourFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid case", "input", "output", false},
        {"error case", "bad", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Troubleshooting

### Tests Requiring Root
Some tests skip automatically if not running as root:
```
--- SKIP: TestRunInstall/with_valid_aptfile_as_root (0.00s)
    install_test.go:105: Skipping test - requires root privileges
```

This is expected and normal.

### System Command Dependencies
Tests that interact with `apt-get`, `dpkg-query`, or `add-apt-repository` will:
- Skip if the command is not available
- Handle permission errors gracefully
- Not actually modify the system

## Next Steps

1. Run `make test-coverage-html` to see detailed coverage
2. Add tests for new features before implementing them
3. Maintain coverage above 80%
4. Run `make test` before committing

For more details, see [TEST_COVERAGE.md](TEST_COVERAGE.md).

