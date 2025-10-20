package main

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/nats-io/nats.go"
	"github.com/telemac/goutils/task"
	"github.com/telemac/nats_service"
	"log"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Create the greeting service
	svcConfig := nats_service.ServiceConfig{
		Nc:          nc,
		Name:        "greeter-service",
		Version:     "1.0.0",
		Description: "A simple greeting service",
		Metadata: map[string]string{
			"author": "nats-service-example",
		},
	}

	greeterService, err := nats_service.NewService(svcConfig)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	defer greeterService.Stop()

	type greetRequest struct {
		Name string `json:"name"`
	}

	type greetResponse struct {
		Greeting string `json:"greeting"`
		Headers  any    `json:"headers"`
	}

	greetEndpoint := nats_service.EndpointConfig{
		Name: "greet",
		Handler: func(req any, msg *nats.Msg, service *nats_service.Service) (any, error) {
			var gr greetRequest
			a := msg.Header.Values("a")
			_ = a
			err := mapstructure.Decode(req, &gr)
			if err != nil {
				return nil, err
			}
			return greetResponse{
				Greeting: fmt.Sprintf("Hello, %s!", gr.Name),
				Headers:  msg.Header,
			}, nil
		},
		RequestTimeout: 0,
		Metadata:       nil,
		QueueGroup:     "",
		Subject:        "test",
	}
	err = greeterService.AddEndpoint(greetEndpoint)
	if err != nil {
		log.Fatalf("Failed to add greet: %v", err)
	}

	<-ctx.Done()

}
