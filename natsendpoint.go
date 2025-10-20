package nats_service

import "github.com/nats-io/nats.go/micro"

var _ Endpointer = (*NatsEndpoint)(nil)

type NatsEndpoint struct {
	config EndpointConfig
}

func (e *NatsEndpoint) GetConfig() EndpointConfig {
	return e.config
}

func (e *NatsEndpoint) SetConfig(config EndpointConfig) {
	e.config = config
}

func (e *NatsEndpoint) Handle(micro.Request) {
	panic("implement me")
}
