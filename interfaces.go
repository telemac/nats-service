package nats_service

import (
	"github.com/nats-io/nats.go/micro"
)

type Servicer interface {
	Stop() error
	GetConfig() ServiceConfig
	AddEndpoint(config EndpointConfig) error
}

type Endpointer interface {
	micro.Handler
	GetConfig() EndpointConfig
	SetConfig(config EndpointConfig)
}
