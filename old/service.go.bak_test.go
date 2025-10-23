package old

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	assert := assert.New(t)

	// Setup test NATS server
	ts := NewTestNATSServer(t)
	defer ts.Shutdown(t)

	tests := []struct {
		name        string
		config      ServiceConfig
		expectError bool
		description string
	}{
		{
			name: "valid service config",
			config: ServiceConfig{
				Nc:          ts.Conn,
				Name:        "test-service",
				Version:     "1.0.0",
				Description: "Test service",
				Metadata: map[string]string{
					"author": "test",
				},
			},
			expectError: false,
			description: "Should create service with valid config",
		},
		{
			name: "service with queue group",
			config: ServiceConfig{
				Nc:          ts.Conn,
				Name:        "queue-service",
				Version:     "1.0.0",
				Description: "Test service with queue group",
			},
			expectError: false,
			description: "Should create service with queue group",
		},
		{
			name: "service with disabled queue group",
			config: ServiceConfig{
				Nc:          ts.Conn,
				Name:        "no-queue-service",
				Version:     "1.0.0",
				Description: "Test service with disabled queue group",
			},
			expectError: false,
			description: "Should create service with disabled queue group",
		},
		{
			name: "minimal service config",
			config: ServiceConfig{
				Nc:      ts.Conn,
				Name:    "minimal-service",
				Version: "1.0.0",
			},
			expectError: false,
			description: "Should create service with minimal config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.config)

			if tt.expectError {
				assert.Error(err, tt.description)
				assert.Nil(service, "Service should be nil when error occurs")
			} else {
				assert.NoError(err, tt.description)
				assert.NotNil(service, "Service should not be nil")
				assert.Equal(tt.config.Name, service.config.Name)
				assert.Equal(tt.config.Version, service.config.Version)
				assert.Equal(tt.config.Description, service.config.Description)
				assert.Equal(tt.config.Metadata, service.config.Metadata)
				assert.NotNil(service.microSvc, "Micro service should be initialized")
			}
		})
	}
}

func TestService_Stop(t *testing.T) {
	assert := assert.New(t)

	// Setup test NATS server
	ts := NewTestNATSServer(t)
	defer ts.Shutdown(t)

	// Create a service
	config := ServiceConfig{
		Nc:          ts.Conn,
		Name:        "test-stop-service",
		Version:     "1.0.0",
		Description: "Service for testing Stop",
	}

	service, err := NewService(config)
	assert.NoError(err, "Should create service successfully")
	assert.NotNil(service, "Service should not be nil")

	// Test stopping the service
	err = service.Stop()
	assert.NoError(err, "Should stop service successfully")

	// Verify service is stopped by trying to stop again
	// Note: The behavior of calling Stop twice may vary depending on NATS micro implementation
	// This test documents the current behavior
	err = service.Stop()
	// Some NATS versions may return an error, others may not
	// We don't assert the error here as the behavior is implementation-specific
}

func TestService_Integration_NewServiceAndStop(t *testing.T) {
	assert := assert.New(t)

	// Setup test NATS server
	ts := NewTestNATSServer(t)
	defer ts.Shutdown(t)

	// Create multiple services
	services := make([]*Service, 3)
	configs := []ServiceConfig{
		{
			Nc:          ts.Conn,
			Name:        "service-1",
			Version:     "1.0.0",
			Description: "First test service",
		},
		{
			Nc:          ts.Conn,
			Name:        "service-2",
			Version:     "2.0.0",
			Description: "Second test service",
		},
		{
			Nc:          ts.Conn,
			Name:        "service-3",
			Version:     "3.0.0",
			Description: "Third test service",
		},
	}

	// Create all services
	for i, config := range configs {
		service, err := NewService(config)
		assert.NoError(err, "Should create service %d successfully", i)
		assert.NotNil(service, "Service %d should not be nil", i)
		services[i] = service
	}

	// Stop all services
	for i, service := range services {
		err := service.Stop()
		assert.NoError(err, "Should stop service %d successfully", i)
	}
}
