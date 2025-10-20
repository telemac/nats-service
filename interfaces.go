package nats_service

import (
	"github.com/nats-io/nats.go/micro"
	"log/slog"
)

type Servicer interface {
	Stop() error
	GetServiceConfig() ServiceConfig
	AddEndpoint(config EndpointConfig) error
	Logger() *slog.Logger
}

type Endpointer interface {
	micro.Handler
	Name() string
	Service() Servicer
	GetEndpointConfig() EndpointConfig
	SetEndpointConfig(config EndpointConfig)
}
