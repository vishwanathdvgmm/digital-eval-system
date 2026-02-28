package core

import (
	"sync"
)

// ServiceRegistry is a simple in-process registry for internal services.
// Later this can be extended to support discovery or health checks.
type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]interface{}
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]interface{}),
	}
}

func (sr *ServiceRegistry) Register(name string, svc interface{}) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.services[name] = svc
}

func (sr *ServiceRegistry) Get(name string) (interface{}, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	v, ok := sr.services[name]
	return v, ok
}

func (sr *ServiceRegistry) MustGet(name string) interface{} {
	v, ok := sr.Get(name)
	if !ok {
		panic("service not found: " + name)
	}
	return v
}
