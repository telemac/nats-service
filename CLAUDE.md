# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go NATS microservice package that provides a clean, interface-based design for building NATS services. The package is designed to be imported by other projects and follows standard Go package conventions.

## Architecture

The codebase follows an interface-based design with these key components:

- **Core Interfaces** (`interfaces.go`):
  - `Servicer`: Interface for service lifecycle management
  - `Endpointer`: Interface for endpoint handling

- **Service Implementation** (`service.go`):
  - `Service`: Main service implementation that manages NATS microservice configuration
  - Handles endpoint registration, service start/stop, and configuration

- **Configuration Types** (`types.go`):
  - `ServiceConfig`: Service-level configuration including NATS connection, logger, metadata
  - `EndpointConfig`: Endpoint-level configuration with metadata, queue groups, subjects

- **Endpoint Framework** (`endpoint.go`):
  - `Endpoint`: Base endpoint implementation that can be embedded in specific endpoints

## Key Patterns

1. **Interface-based design**: Services implement `Servicer`, endpoints implement `Endpointer`
2. **Embedded endpoints**: Endpoints embed `natsservice.Endpoint` and implement `Handle(micro.Request)`
3. **Structured logging**: Uses `log/slog` for logging throughout the codebase
4. **Package at root**: Main package files are at the root level for easy importing

## Development Commands

### Build
```bash
go build ./...
```

### Run Example Service
```bash
go run ./cmd/basic-service
```

### Run Tests
```bash
go test ./...
```

### Run Specific Test
```bash
go test -v ./path/to/package
```

### Lint/Vet
```bash
go vet ./...
```

### Format Code
```bash
go fmt ./...
```

## Project Structure

```
nats-service/
├── service.go              # Main service implementation
├── interfaces.go           # Public interfaces
├── types.go               # Configuration types
├── endpoint.go            # Base endpoint implementation
├── pkg/                   # Supporting packages
│   └── counter/           # Counter utility package
├── examples/              # Usage examples
│   └── basic/
│       ├── service.go     # Example service implementation
│       └── endpoints/
│           ├── add.go
│           └── say.go
├── cmd/                   # Executable examples
│   └── basic-service/
│       └── main.go
├── README.md              # Package documentation for consumers
└── CLAUDE.md              # This file - development guide
```

## Testing Standards

When writing Go tests, initialize assertions with `assert := assert.New(t)` at the beginning of test functions and use `assert.Equal()`, `assert.NoError()`, etc. without the `t` parameter.

## Dependencies

- NATS.io for messaging and microservices
- testify for testing assertions
- goutils for utilities

## Package Usage

This is a library package meant to be imported by other projects:
```go
import "github.com/telemac/natsservice"
```