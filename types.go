package nats_service

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"log/slog"
)

type ServiceConfig struct {
	Ctx         context.Context // TODO : handle or remove
	Nc          *nats.Conn
	Logger      *slog.Logger
	Name        string `json:"name"`
	Prefix      string
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

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

type EndpointConfig struct {
	Service  Servicer
	Name     string     `json:"name"`
	Endpoint Endpointer `json:"-"` // Non-serializable
	//RequestTimeout time.Duration     `json:"request_timeout"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	QueueGroup string            `json:"queue_group,omitempty"`
	Subject    string            `json:"subject,omitempty"`
}

func (ec *EndpointConfig) Validate() error {
	if ec.Service == nil {
		return errors.New("service required")
	}
	if ec.Name == "" {
		return errors.New("service name required")
	}
	if ec.Endpoint == nil {
		return errors.New("endpoint required")
	}

	return nil
}
