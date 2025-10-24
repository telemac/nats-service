package natsservice

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"log/slog"
)

// ServiceConfig holds configuration for NATS services
type ServiceConfig struct {
	Ctx         context.Context   // Service context for cancellation
	Nc          *nats.Conn        // NATS connection
	Logger      *slog.Logger      // Service logger
	Name        string            `json:"name"` // Service name
	Prefix      string            // Subject prefix for endpoints
	Version     string            `json:"version"`            // Service version
	Description string            `json:"description"`        // Service description
	Metadata    map[string]string `json:"metadata,omitempty"` // Additional metadata
}

// Validate checks that all required fields are present
func (sc *ServiceConfig) Validate() error {
	if sc.Ctx == nil {
		return errors.New("missing context")
	}
	if sc.Nc == nil {
		return errors.New("nats connection required")
	}
	if sc.Logger == nil {
		return errors.New("logger required")
	}
	if sc.Name == "" {
		return errors.New("service name required")
	}
	return nil
}

// EndpointConfig holds configuration for individual endpoints
type EndpointConfig struct {
	Service    Servicer          // Associated service
	Name       string            `json:"name"`                  // Endpoint name
	Endpoint   Endpointer        `json:"-"`                     // Non-serializable endpoint implementation
	Metadata   map[string]string `json:"metadata,omitempty"`    // Endpoint metadata
	QueueGroup string            `json:"queue_group,omitempty"` // Queue group name
	Subject    string            `json:"subject,omitempty"`     // Custom subject
}

// Validate checks that all required fields are present
func (ec *EndpointConfig) Validate() error {
	if ec.Service == nil {
		return errors.New("service required")
	}
	if ec.Name == "" {
		return errors.New("endpoint name required")
	}
	if ec.Endpoint == nil {
		return errors.New("endpoint required")
	}
	return nil
}
