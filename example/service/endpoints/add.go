package endpoints

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/nats_service"
)

type Add struct {
	nats_service.NatsEndpoint
	Count int
}

func (e *Add) Handle(req micro.Request) {
	log := e.Service().Logger()
	e.Count++
	log.Info("add handler called", "count", e.Count)
	req.RespondJSON(e.Count)
}
