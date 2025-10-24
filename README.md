# nats-service

A Go package that provides a clean, interface-based abstraction layer over NATS micro for building microservices. This library focuses on simplicity and type safety while working with existing NATS connections.

## Features

- **Interface-based design**: Clean separation between service and endpoint implementations
- **Type-safe configuration**: Structured configuration with validation
- **Embedded endpoint pattern**: Easy to create endpoints by embedding the base endpoint type
- **No connection management**: Works with your existing NATS connections
- **Structured logging**: Built-in slog integration
- **Graceful shutdown**: Proper service lifecycle management

## Installation

```bash
go get github.com/telemac/nats_service
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/nats-io/nats.go"
    "github.com/telemac/nats_service"
    "log/slog"
)

func main() {
    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Connect to NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        panic(err)
    }
    defer nc.Close()

    // Create service
    var svc nats_service.Service
    err = svc.Start(&nats_service.ServiceConfig{
        Ctx:         ctx,
        Nc:          nc,
        Logger:      slog.Default(),
        Name:        "my-service",
        Version:     "1.0.0",
        Description: "My awesome service",
    })
    if err != nil {
        panic(err)
    }
    defer svc.Stop()

    // Add endpoints
    err = svc.AddEndpoint(&nats_service.EndpointConfig{
        Endpoint: &MyEndpoint{},
    })
    if err != nil {
        panic(err)
    }

    // Service is now running
    select {}
}
```

## Core Concepts

### Service

The main component that manages NATS microservice configuration and lifecycle.

```go
type Servicer interface {
    Start(*ServiceConfig) error
    Stop() error
    GetServiceConfig() *ServiceConfig
    AddEndpoint(config *EndpointConfig) error
    AddEndpoints(configs ...*EndpointConfig) error
    Logger() *slog.Logger
}
```

### Endpoint

Endpoints handle incoming NATS requests. Embed the base `Endpoint` type and implement the `Handle` method.

```go
type MyEndpoint struct {
    nats_service.Endpoint
}

func (e *MyEndpoint) Name() string {
    return "my-endpoint"
}

func (e *MyEndpoint) Handle(req micro.Request) {
    // Handle the request
    req.RespondJSON(map[string]string{"message": "Hello, World!"})
}
```

### Configuration

Service and endpoint configuration is done through structured types:

```go
type ServiceConfig struct {
    Ctx         context.Context
    Nc          *nats.Conn
    Logger      *slog.Logger
    Name        string
    Prefix      string
    Version     string
    Description string
    Metadata    map[string]string
}

type EndpointConfig struct {
    Service    Servicer
    Name       string
    Endpoint   Endpointer
    Metadata   map[string]string
    QueueGroup string
    Subject    string
}

```

Note: When registering endpoints, the service merges metadata from both the EndpointConfig.Metadata field
and the endpoint's Metadata() method. If both are present, config metadata takes precedence.

## Examples

See the [examples](./examples) directory for complete working examples:

- **[Basic Service](./examples/basic/)**: A simple service with multiple endpoints
- **[Command-line Tool](./cmd/basic-service/)**: Executable version of the basic example

## Development Commands

```bash
# Build the package
go build ./...

# Run the example service
go run ./cmd/basic-service

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet for potential issues
go vet ./...
```

## Project Structure

```
nats-service/
├── service.go          # Main service implementation
├── interfaces.go       # Public interfaces
├── types.go           # Configuration types
├── endpoint.go        # Base endpoint implementation
├── pkg/               # Supporting packages
│   └── counter/       # Counter utility
├── examples/          # Usage examples
│   └── basic/
├── cmd/               # Executable examples
│   └── basic-service/
└── README.md
```

## Dependencies

- Go 1.24+
- github.com/nats-io/nats.go
- github.com/telemac/goutils

## License

[Your License Here]