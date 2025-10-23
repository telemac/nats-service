package nats_service

import "github.com/nats-io/nats.go/micro"

// Ensure Endpoint implements Endpointer interface
var _ Endpointer = (*Endpoint)(nil)

// Endpoint is a base implementation of the Endpointer interface
type Endpoint struct {
	config *EndpointConfig // Endpoint configuration
}

// GetEndpointConfig returns the endpoint configuration
func (e *Endpoint) GetEndpointConfig() *EndpointConfig {
	return e.config
}

// SetEndpointConfig sets the endpoint configuration
func (e *Endpoint) SetEndpointConfig(config *EndpointConfig) {
	e.config = config
}

// Service returns the associated service
func (e *Endpoint) Service() Servicer {
	return e.config.Service
}

// Handle processes incoming requests - must be implemented by embedding types
func (e *Endpoint) Handle(micro.Request) {
	panic("implement me")
}

// Name returns the endpoint name - must be implemented by embedding types
func (e *Endpoint) Name() string {
	return ""
}

// Metadata returns endpoint metadata - must be implemented by embedding types
func (e *Endpoint) Metadata() map[string]string {
	return nil
}
