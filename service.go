package nats_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"time"
)

// check that Service satisfies ServiceInterface
//var _ Servicer = (*Service)(nil)

type Service struct {
	config   ServiceConfig
	microSvc micro.Service
}

type ServiceConfig struct {
	Ctx                context.Context // TOCO : handle or remove
	Nc                 *nats.Conn
	Name               string `json:"name"`
	prefix             string
	Version            string            `json:"version"`
	Description        string            `json:"description"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	QueueGroup         string            `json:"queue_group"`
	QueueGroupDisabled bool              `json:"queue_group_disabled"`
}

func (sc *ServiceConfig) Validate() error {
	if sc.Nc == nil {
		return errors.New("nats connection required")
	}
	if sc.Name == "" {
		return errors.New("service name required")
	}
	return nil
}

// NewService initializes and returns a new Service instance based on the provided ServiceConfig.
func NewService(serviceConfig ServiceConfig) (*Service, error) {
	err := serviceConfig.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid service config: %w", err)
	}
	service := &Service{
		config: serviceConfig,
	}
	// build micro.AddService configuration
	microConfig := micro.Config{
		Name:               serviceConfig.Name,
		Version:            serviceConfig.Version,
		Metadata:           serviceConfig.Metadata,
		QueueGroup:         serviceConfig.QueueGroup,
		QueueGroupDisabled: serviceConfig.QueueGroupDisabled,
	}
	service.microSvc, err = micro.AddService(service.config.Nc, microConfig)
	return service, err
}

// Stop drains the endpoint subscriptions and marks the service as stopped.
func (svc *Service) Stop() error {
	return svc.microSvc.Stop()
}

func (s *Service) Config() ServiceConfig {
	return s.config
}

type EndpointHandler = func(req any, msg *nats.Msg, service *Service) (any, error)

type EndpointConfig struct {
	Name           string            `json:"name"`
	Handler        EndpointHandler   `json:"-"` // Non-serializable
	RequestTimeout time.Duration     `json:"request_timeout"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	QueueGroup     string            `json:"queue_group,omitempty"`
	Subject        string            `json:"subject,omitempty"`
	Service        *Service
}

func (ec *EndpointConfig) Validate() error {
	if ec.Name == "" {
		return errors.New("endpoint name is required")
	}
	if ec.Handler == nil {
		return errors.New("endpoint handler is required")
	}
	return nil
}

type handler struct {
	ec *EndpointConfig
}

func newHandler(ec *EndpointConfig) *handler {
	return &handler{ec: ec}
}

func (h *handler) Handle(request micro.Request) {
	if h.ec.Handler == nil {
		panic("endpoint handler is required")
	}
	// Extract the request data from the NATS message
	var reqData any
	if request.Data() != nil {
		err := json.Unmarshal(request.Data(), &reqData)
		if err != nil {
			if len(request.Data()) > 0 {
				request.Error("500", err.Error(), nil)
				return
			}
		}
	}

	// Create a NATS message from the request data for the handler
	// Note: We're creating a minimal message since micro.Request doesn't expose the full NATS message
	msg := &nats.Msg{
		Subject: request.Subject(),
		Reply:   request.Reply(),
		Data:    request.Data(),
	}
	headers := request.Headers()
	if len(headers) > 0 {
		msg.Header = make(nats.Header)
		for header, values := range headers {
			for _, value := range values {
				msg.Header.Add(header, value)
			}
		}
	}

	// Call the handler and get the response
	response, err := h.ec.Handler(reqData, msg, h.ec.Service)
	if err != nil {
		// Send error response
		if err := request.Error("500", err.Error(), nil); err != nil {
			// Log error if we can't send error response
			// In a real implementation, use proper logging
		}
		return
	}

	// Send successful response
	if response != nil {
		// For now, handle string responses
		if str, ok := response.(string); ok {
			if err := request.Respond([]byte(str)); err != nil {
				// Log error if we can't send response
				// In a real implementation, use proper logging
			}
		} else {
			// Try to marshal as JSON for other types
			if err := request.RespondJSON(response); err != nil {
				// Log error if we can't send JSON response
				// In a real implementation, use proper logging
			}
		}
	} else {
		// Send empty response for nil response
		if err := request.Respond([]byte{}); err != nil {
			// Log error if we can't send empty response
			// In a real implementation, use proper logging
		}
	}
}

func (svc *Service) AddEndpoint(config EndpointConfig) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	config.Service = svc

	var opts []micro.EndpointOpt

	// Subject (obligatoire dans la plupart des cas)
	if config.Subject != "" {
		opts = append(opts, micro.WithEndpointSubject(config.Subject))
	} else {
		opts = append(opts, micro.WithEndpointSubject(svc.config.Name+"."+config.Name))
	}

	// Métadonnées
	if len(config.Metadata) > 0 {
		opts = append(opts, micro.WithEndpointMetadata(config.Metadata))
	}

	// Groupe de queue (activé ou désactivé)
	if config.QueueGroup != "" {
		opts = append(opts, micro.WithEndpointQueueGroup(config.QueueGroup))
	} else {
		opts = append(opts, micro.WithEndpointQueueGroupDisabled())
	}

	return svc.microSvc.AddEndpoint(config.Name, newHandler(&config), opts...)
}
