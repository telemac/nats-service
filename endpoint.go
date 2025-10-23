package nats_service

import "github.com/nats-io/nats.go/micro"

var _ Endpointer = (*Endpoint)(nil)

type Endpoint struct {
	config *EndpointConfig
}

func (e *Endpoint) GetEndpointConfig() *EndpointConfig {
	return e.config
}

func (e *Endpoint) SetEndpointConfig(config *EndpointConfig) {
	e.config = config
}

func (e *Endpoint) Service() Servicer {
	return e.config.Service
}

func (e *Endpoint) Handle(micro.Request) {
	panic("implement me")
}

func (e *Endpoint) Name() string {
	return ""
}

func (e *Endpoint) Metadata() map[string]string {
	return nil
}
