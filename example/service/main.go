package main

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/telemac/goutils/task"
	"github.com/telemac/nats_service"
	"log/slog"
	"time"
)

type Add struct {
	nats_service.NatsEndpoint
	Count int
}

func (e *Add) Handle(micro.Request) {
	e.Count++
}

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	log := slog.Default().With("version", "0.0.1")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		return
	}
	defer nc.Close()

	log.Info("starting nats service")
	service, err := nats_service.StartNatsService(nats_service.ServiceConfig{
		Ctx:         ctx,
		Nc:          nc,
		Name:        "math",
		Version:     "0.0.1",
		Description: "make additions and substractions",
		Prefix:      "a.b.c",
		Metadata: map[string]string{
			"authro": "Alexandre HEIM",
		},
	})
	if err != nil {
		slog.Error("Failed to run NATS service", "error", err)
		return
	}

	service.GetConfig().Metadata["prefix"] = "a.b.c" // how to modify service metadatas

	var addEndpoint = &Add{}

	err = service.AddEndpoint(nats_service.EndpointConfig{
		Service:        nil,
		Name:           "add",
		Endpoint:       addEndpoint,
		RequestTimeout: 0,
		Metadata:       nil,
		QueueGroup:     "",
		Subject:        "",
	})
	if err != nil {
		slog.Error("Failed to add endpoint", "error", err)
		return
	}

	<-ctx.Done()

	log.Info("shutting down NATS service")
	err = service.Stop()
	if err != nil {
		slog.Error("Failed to stop NATS service", "error", err)
	}

}
