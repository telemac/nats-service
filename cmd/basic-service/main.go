package main

import (
	"github.com/nats-io/nats.go"
	"github.com/telemac/goutils/task"
	"github.com/telemac/nats-service"
	"github.com/telemac/nats-service/examples/basic"
	"github.com/telemac/nats-service/examples/basic/endpoints"
	"log/slog"
	"time"
)

func main() {
	// Create cancellable context with 5s timeout
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	// Initialize logger with version
	log := slog.Default().With("version", "0.0.1")

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		return
	}
	defer nc.Close()

	// Create and start service
	log.Info("starting nats service")
	var basicService basic.BasicService
	err = basicService.Start(&nats_service.ServiceConfig{
		Ctx:         ctx,
		Nc:          nc,
		Logger:      log,
		Name:        "basic-service",
		Version:     "0.0.1",
		Description: "service example",
		Metadata: map[string]string{
			"author": "Alexandre HEIM",
		},
	})
	if err != nil {
		slog.Error("Failed to run NATS service", "error", err)
		return
	}

	// Ensure graceful shutdown
	defer func() {
		log.Info("shutting down NATS service")
		err = basicService.Stop()
		if err != nil {
			slog.Error("Failed to stop NATS service", "error", err)
		}
	}()

	// Add metadata at runtime
	basicService.GetServiceConfig().Metadata["prefix"] = "a.b.c"

	// Register endpoints
	err = basicService.AddEndpoints(
		&nats_service.EndpointConfig{
			Endpoint: &endpoints.Add{},
		},
		&nats_service.EndpointConfig{
			Endpoint: &endpoints.Say{},
		},
	)
	if err != nil {
		slog.Error("Failed to add endpoint", "error", err, "endpoint", &endpoints.Add{})
		return
	}

	// Wait for context cancellation
	<-ctx.Done()
}
