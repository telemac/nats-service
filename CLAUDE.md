# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go library that provides a clean, type-safe abstraction layer over NATS micro for experienced developers who want simpler abstractions. The library works with existing NATS connections and focuses on Publish/Subscribe and Request/Reply patterns without JetStream complexity.

## Key Architecture

### Core Components

- **Service**: Main component that wraps NATS micro functionality, always created with a Config
- **Publisher**: Handles message publishing to NATS subjects
- **Subscriber**: Manages message subscriptions with optional queue groups
- **Requester**: Sends request messages and handles responses
- **Responder**: Handles incoming requests and sends replies

### Design Philosophy

- **KISS Principle**: Simple, minimal abstractions over NATS micro
- **Type Safety**: Full generics support for compile-time checking
- **Existing Connection**: Works with user-provided NATS connections, doesn't manage connections
- **Developer Experience**: Fluent APIs and clean patterns
- **Zero Overhead**: Thin wrapper layer with minimal performance impact

## Development Commands

```bash
# Build the library
go build ./...

# Run example
go run ./example

# Run tests (when they exist)
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestFunction ./pkg/nats_service

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run

# Run benchmarks
go test -bench=. ./...
```

## Common Patterns

### Service Creation
All services are created with a Config structure:

```go
type Config struct {
    Name     string
    Logger   *slog.Logger
	...
}

func NewService(cfg Config) (*Service, error) {
    // Initialize with config values
}
```

## Module Structure

- Module name: `github.com/telemac/nats_service`
- Go version: >=1.24
- Main dependency: `github.com/nats-io/nats.go`
- Package structure: `pkg/nats_service` for library code, `example` for usage examples
- Logging: Uses `log/slog` for structured logging

## Examples

The example directory should contain:
- `publisher.go`: Simple message publisher
- `subscriber.go`: Simple message subscriber with optional queue group
- `greeter.go`: Request/Reply greeter service with optional queue group

## Testing

When adding tests:
- Use testify/assert for assertions
- Initialize with `assert := assert.New(t)`
- Place test files alongside source files with `_test.go` suffix
- Mock NATS connections for unit tests
- Test both queue group and non-queue group scenarios