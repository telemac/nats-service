package nats_service

import (
	"fmt"
	"github.com/nats-io/nats.go/micro"
	"log/slog"
	"sync"
)

var _ Servicer = (*Service)(nil)

type Service struct {
	config   *ServiceConfig
	mu       sync.RWMutex
	microSvc micro.Service
}

//func StartService(config *ServiceConfig) (*Service, error) {
//	var svc Service
//	err := svc.Start(config)
//	return &svc, err
//}

func (svc *Service) Start(config *ServiceConfig) error {
	err := config.Validate()
	if err != nil {
		return fmt.Errorf("invalid service config: %w", err)
	}
	svc.config = config
	// build micro.AddService configuration
	microConfig := micro.Config{
		Name:        svc.config.Name,
		Version:     svc.config.Version,
		Description: svc.config.Description,
		Metadata:    svc.config.Metadata,
	}
	svc.microSvc, err = micro.AddService(svc.config.Nc, microConfig)
	if err != nil {
		return err
	}
	if svc.config.Prefix != "" {
		svc.microSvc.AddGroup(svc.config.Prefix)
	}
	return err
}

func (svc *Service) Stop() error {
	return svc.microSvc.Stop()
}

func (svc *Service) GetServiceConfig() ServiceConfig {
	return *svc.config
}

func (svc *Service) AddEndpoint(config EndpointConfig) error {
	if config.Name == "" {
		config.Name = config.Endpoint.Name()
	}
	err := svc.config.Validate()
	if err != nil {
		return fmt.Errorf("invalid endpoint config: %w", err)
	}
	config.Service = svc
	config.Endpoint.SetEndpointConfig(config)

	var opts []micro.EndpointOpt

	// Subject (obligatoire dans la plupart des cas)
	if config.Subject != "" {
		opts = append(opts, micro.WithEndpointSubject(config.Subject))
	} else {
		if svc.config.Prefix != "" {
			opts = append(opts, micro.WithEndpointSubject(svc.config.Prefix+"."+svc.config.Name+"."+config.Name))
		} else {
			opts = append(opts, micro.WithEndpointSubject(svc.config.Name+"."+config.Name))
		}
	}

	// Métadonnées
	if len(config.Metadata) > 0 && len(config.Endpoint.Metadata()) == 0 {
		opts = append(opts, micro.WithEndpointMetadata(config.Metadata))
	} else if len(config.Metadata) == 0 && len(config.Endpoint.Metadata()) > 0 {
		opts = append(opts, micro.WithEndpointMetadata(config.Metadata))
	} else {
		// merge metas
		metas := make(map[string]string)
		for k, v := range config.Endpoint.Metadata() {
			metas[k] = v
		}
		for k, v := range config.Metadata {
			metas[k] = v
		}
		opts = append(opts, micro.WithEndpointMetadata(metas))
	}

	// Groupe de queue (activé ou désactivé)
	if config.QueueGroup != "" {
		opts = append(opts, micro.WithEndpointQueueGroup(config.QueueGroup))
	} else {
		opts = append(opts, micro.WithEndpointQueueGroupDisabled())
	}

	return svc.microSvc.AddEndpoint(config.Name, config.Endpoint, opts...)
}

func (svc *Service) Logger() *slog.Logger {
	return svc.config.Logger
}
