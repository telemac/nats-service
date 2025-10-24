package basic

import (
	"github.com/telemac/nats-service"
	"github.com/telemac/nats-service/pkg/counter"
)

// Ensure BasicService implements Servicer interface
var _ nats_service.Servicer = (*BasicService)(nil)

// BasicService extends nats_service.Service with a counter
type BasicService struct {
	nats_service.Service
	Counter counter.Counter // Service-level counter for tracking requests
}
