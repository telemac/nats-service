package nats_service

import (
	"context"
	"github.com/nats-io/nats.go/micro"
)

type Servicer interface {
	Run(ctx context.Context, config ServiceConfig) error
	Stop() error
	AddEndpoint(endPoint Endpointer)
}

type Endpointer interface {
	micro.Handler
	Service() Servicer
	GetConfig() EndpointConfig
}
