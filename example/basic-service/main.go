package main

import (
	"github.com/nats-io/nats.go"
	"github.com/samber/do/v2"
	"github.com/telemac/goutils/task"
	"github.com/telemac/nats_service"
	"github.com/telemac/nats_service/example/basic-service/basicservice"
	"github.com/telemac/nats_service/example/basic-service/endpoints"
	"github.com/telemac/nats_service/pkg/counter"
	"log/slog"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	injector := do.New()
	do.Provide(injector, func(i do.Injector) (*counter.Counter, error) {
		return &counter.Counter{}, nil
	})

	//counter := do.MustInvoke[*counter.Counter](injector)
	//counter.Increment(1)

	log := slog.Default().With("version", "0.0.1")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		return
	}
	defer nc.Close()

	log.Info("starting nats service")
	var basicService basicservice.BasicService
	err = basicService.Start(&nats_service.ServiceConfig{
		Ctx:         ctx,
		Nc:          nc,
		Logger:      log,
		Name:        "basic-service",
		Version:     "0.0.1",
		Description: "service example",
		//Prefix:      "a.b.c",
		Metadata: map[string]string{
			"authro": "Alexandre HEIM",
		},
	})
	//basicService, err := nats_service.StartService(&nats_service.ServiceConfig{
	//	Ctx:         ctx,
	//	Nc:          nc,
	//	Logger:      log,
	//	Name:        "basic-service",
	//	Version:     "0.0.1",
	//	Description: "service example",
	//	//Prefix:      "a.b.c",
	//	Metadata: map[string]string{
	//		"authro": "Alexandre HEIM",
	//	},
	//})
	if err != nil {
		slog.Error("Failed to run NATS service", "error", err)
		return
	}
	defer func() {
		log.Info("shutting down NATS service")
		err = basicService.Stop()
		if err != nil {
			slog.Error("Failed to stop NATS service", "error", err)
		}
	}()

	basicService.GetServiceConfig().Metadata["prefix"] = "a.b.c" // how to modify service metadatas

	endpoints := []nats_service.Endpointer{
		&endpoints.Add{},
		&endpoints.Say{},
	}

	for _, endpoint := range endpoints {
		err = basicService.AddEndpoint(nats_service.EndpointConfig{
			Endpoint: endpoint,
		})
		if err != nil {
			slog.Error("Failed to add endpoint", "error", err, "endpoint", endpoint)
			return
		}
	}

	<-ctx.Done()

}
