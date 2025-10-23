# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go NATS microservice framework that provides an interface-based design for building NATS services. The project uses NATS.io for messaging and provides a clean abstraction layer for creating microservices with endpoints.

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
2. **Embedded endpoints**: Endpoints embed `nats_service.Endpoint` and implement `Handle(micro.Request)`
4. **Structured logging**: Uses `log/slog` for logging throughout the codebase

## Development Commands

### Build
```bash
go build ./...
```

### Run Example Service
```bash
go run ./example/basic-service/main.go
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

- `example/basic-service/`: Complete example service implementation
  - `basicservice/`: Custom service type embedding the base Service
  - `endpoints/`: Example endpoints (Add, Say)
- `pkg/counter/`: Shared counter utility package
- `old/`: Backup/legacy test files

## Testing Standards

When writing Go tests, initialize assertions with `assert := assert.New(t)` at the beginning of test functions and use `assert.Equal()`, `assert.NoError()`, etc. without the `t` parameter.

## Dependencies

- NATS.io for messaging and microservices
- testify for testing assertions
- goutils for utilities