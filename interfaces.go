package natsservice

import (
	"github.com/nats-io/nats.go/micro"
	"log/slog"
)

// Servicer defines the lifecycle and management interface for NATS services
type Servicer interface {
	Start(*ServiceConfig) error                    // Start the service with configuration
	Stop() error                                   // Stop the service gracefully
	GetServiceConfig() *ServiceConfig              // Get current service configuration
	AddEndpoint(config *EndpointConfig) error      // Add a single endpoint
	AddEndpoints(configs ...*EndpointConfig) error // Add multiple endpoints
	Logger() *slog.Logger                          // Get the service logger
}

// Endpointer defines the interface for service endpoints
type Endpointer interface {
	micro.Handler                             // NATS micro handler interface
	Name() string                             // Get endpoint name
	Metadata() map[string]string              // Get endpoint metadata
	Service() Servicer                        // Get associated service
	GetEndpointConfig() *EndpointConfig       // Get endpoint configuration
	SetEndpointConfig(config *EndpointConfig) // Set endpoint configuration
}
