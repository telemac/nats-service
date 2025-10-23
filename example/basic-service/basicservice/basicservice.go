package basicservice

import (
	"github.com/telemac/nats_service"
	"github.com/telemac/nats_service/pkg/counter"
)

var _ nats_service.Servicer = (*BasicService)(nil)

type BasicService struct {
	nats_service.Service
	Counter counter.Counter
}
