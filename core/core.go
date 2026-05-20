// Package core provides the foundational components of v2ray-core.
// It defines the Instance type which represents a running V2Ray instance,
// managing inbound/outbound handlers and routing.
package core

import (
	"context"
	"sync"

	"github.com/v2fly/v2ray-core/v5/common"
)

// Version is the current version of v2ray-core.
const Version = "5.0.0"

// Instance combines all features needed for a working V2Ray instance.
type Instance struct {
	access   sync.Mutex
	features []Feature
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// Feature is the interface for all V2Ray features/components.
type Feature interface {
	common.Runnable
	Type() interface{}
}

// New creates a new Instance based on the given config.
// It initializes all features defined in the configuration.
func New(config *Config) (*Instance, error) {
	ctx, cancel := context.WithCancel(context.Background())
	inst := &Instance{
		ctx:    ctx,
		cancel: cancel,
	}

	if err := inst.applyConfig(config); err != nil {
		cancel()
		return nil, err
	}

	return inst, nil
}

// applyConfig initializes the instance based on the provided Config.
func (s *Instance) applyConfig(config *Config) error {
	if config == nil {
		return newError("config is nil")
	}
	// Additional config application logic will go here
	// as more components are added (inbounds, outbounds, routing, etc.)
	return nil
}

// AddFeature registers a new feature to the instance.
// If the instance is already running, the feature will be started immediately.
func (s *Instance) AddFeature(feature Feature) error {
	s.access.Lock()
	defer s.access.Unlock()

	s.features = append(s.features, feature)

	if s.running {
		if err := feature.Start(); err != nil {
			return newError("failed to start feature").Base(err)
		}
	}
	return nil
}

// GetFeature returns the first registered feature that matches the given type.
func (s *Instance) GetFeature(featureType interface{}) Feature {
	s.access.Lock()
	defer s.access.Unlock()

	for _, f := range s.features {
		if f.Type() == featureType {
			return f
		}
	}
	return nil
}

// Start starts the V2Ray instance, initializing and running all registered features.
func (s *Instance) Start() error {
	s.access.Lock()
	defer s.access.Unlock()

	if s.running {
		return newError("instance is already running")
	}

	for _, f := range s.features {
		if err := f.Start(); err != nil {
			return newError("failed to start feature").Base(err)
		}
	}

	s.running = true
	return nil
}

// Close shuts down the V2Ray instance and all its features.
func (s *Instance) Close() error {
	s.access.Lock()
	defer s.access.Unlock()

	if !s.running {
		return newError("instance is not running")
	}

	s.running = false
	s.cancel()

	var errs []error
	for _, f := range s.features {
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return newError("errors occurred while closing instance").Base(errs[0])
	}
	return nil
}

// Context returns the context associated with this instance.
func (s *Instance) Context() context.Context {
	return s.ctx
}

// IsRunning returns whether the instance is currently active.
func (s *Instance) IsRunning() bool {
	s.access.Lock()
	defer s.access.Unlock()
	return s.running
}
