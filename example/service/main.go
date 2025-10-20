package main

import (
	"github.com/nats-io/nats.go"
	"github.com/telemac/goutils/task"
	"github.com/telemac/nats_service"
	"github.com/telemac/nats_service/example/service/endpoints"
	"log/slog"
	"time"
)

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
		Logger:      log,
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

	service.GetServiceConfig().Metadata["prefix"] = "a.b.c" // how to modify service metadatas

	var addEndpoint = &endpoints.Add{}

	err = service.AddEndpoint(nats_service.EndpointConfig{
		//Name:           "add",
		Endpoint:   addEndpoint,
		Metadata:   map[string]string{"prefix": "a.b.c"},
		QueueGroup: "",
		Subject:    "",
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
