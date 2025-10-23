package nats_service

import (
	"fmt"
	"github.com/nats-io/nats.go/micro"
	"log/slog"
	"sync"
)

// Ensure Service implements Servicer interface
var _ Servicer = (*Service)(nil)

// Service is the default implementation of the Servicer interface
type Service struct {
	config   *ServiceConfig  // Service configuration
	mu       sync.RWMutex    // Mutex for thread safety
	microSvc micro.Service   // NATS micro service instance
}

// Start initializes and starts the NATS microservice
func (svc *Service) Start(config *ServiceConfig) error {
	// Validate configuration
	err := config.Validate()
	if err != nil {
		return fmt.Errorf("invalid service config: %w", err)
	}
	svc.config = config

	// Build micro service configuration
	microConfig := micro.Config{
		Name:        svc.config.Name,
		Version:     svc.config.Version,
		Description: svc.config.Description,
		Metadata:    svc.config.Metadata,
	}

	// Create micro service
	svc.microSvc, err = micro.AddService(svc.config.Nc, microConfig)
	if err != nil {
		return err
	}

	// Add prefix group if specified
	if svc.config.Prefix != "" {
		svc.microSvc.AddGroup(svc.config.Prefix)
	}

	return err
}

// Stop gracefully stops the NATS microservice
func (svc *Service) Stop() error {
	return svc.microSvc.Stop()
}

// GetServiceConfig returns the current service configuration
func (svc *Service) GetServiceConfig() *ServiceConfig {
	return svc.config
}

// AddEndpoint registers a single endpoint with the service
func (svc *Service) AddEndpoint(config *EndpointConfig) error {
	if config == nil {
		return fmt.Errorf("invalid endpoint config")
	}

	// Set endpoint name if not provided
	if config.Name == "" {
		config.Name = config.Endpoint.Name()
	}

	// Validate service configuration
	err := svc.config.Validate()
	if err != nil {
		return fmt.Errorf("invalid endpoint config: %w", err)
	}

	// Link endpoint to service and configuration
	config.Service = svc
	config.Endpoint.SetEndpointConfig(config)

	// Build endpoint options
	var opts []micro.EndpointOpt

	// Configure subject
	if config.Subject != "" {
		opts = append(opts, micro.WithEndpointSubject(config.Subject))
	} else {
		// Generate default subject
		if svc.config.Prefix != "" {
			opts = append(opts, micro.WithEndpointSubject(svc.config.Prefix+"."+svc.config.Name+"."+config.Name))
		} else {
			opts = append(opts, micro.WithEndpointSubject(svc.config.Name+"."+config.Name))
		}
	}

	// Configure metadata
	if len(config.Metadata) > 0 && len(config.Endpoint.Metadata()) == 0 {
		opts = append(opts, micro.WithEndpointMetadata(config.Metadata))
	} else if len(config.Metadata) == 0 && len(config.Endpoint.Metadata()) > 0 {
		opts = append(opts, micro.WithEndpointMetadata(config.Endpoint.Metadata()))
	} else {
		// Merge metadata from both config and endpoint
		metas := make(map[string]string)
		for k, v := range config.Endpoint.Metadata() {
			metas[k] = v
		}
		for k, v := range config.Metadata {
			metas[k] = v
		}
		opts = append(opts, micro.WithEndpointMetadata(metas))
	}

	// Configure queue group
	if config.QueueGroup != "" {
		opts = append(opts, micro.WithEndpointQueueGroup(config.QueueGroup))
	} else {
		opts = append(opts, micro.WithEndpointQueueGroupDisabled())
	}

	// Register endpoint with NATS micro
	return svc.microSvc.AddEndpoint(config.Name, config.Endpoint, opts...)
}

// AddEndpoints registers multiple endpoints with the service
func (svc *Service) AddEndpoints(configs ...*EndpointConfig) error {
	for _, config := range configs {
		err := svc.AddEndpoint(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// Logger returns the service logger
func (svc *Service) Logger() *slog.Logger {
	return svc.config.Logger
}
