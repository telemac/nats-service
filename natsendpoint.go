package nats_service

import "github.com/nats-io/nats.go/micro"

var _ Endpointer = (*NatsEndpoint)(nil)

type NatsEndpoint struct {
	config EndpointConfig
}

func (e *NatsEndpoint) GetEndpointConfig() EndpointConfig {
	return e.config
}

func (e *NatsEndpoint) SetEndpointConfig(config EndpointConfig) {
	e.config = config
}

func (e *NatsEndpoint) Service() Servicer {
	return e.config.Service
}

func (e *NatsEndpoint) Handle(micro.Request) {
	panic("implement me")
}

func (e *NatsEndpoint) Name() string {
	return ""
}

func (e *NatsEndpoint) Metadata() map[string]string {
	return nil
}
