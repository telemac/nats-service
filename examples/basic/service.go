package basic

import (
	"github.com/telemac/natsservice"
	"github.com/telemac/natsservice/pkg/counter"
)

// Ensure BasicService implements Servicer interface
var _ natsservice.Servicer = (*BasicService)(nil)

// BasicService extends natsservice.Service with a counter
type BasicService struct {
	natsservice.Service
	Counter counter.Counter // Service-level counter for tracking requests
}
