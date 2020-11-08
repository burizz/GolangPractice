package serviceregistry

import (
	"fmt"
	"log"
	"reflect"
)

// Service - defines needed methods for each service
type Service interface {
	// Start spawns main process done by the service
	Start()
	// Stop terminates all processes belonging to the service,
	// blocking until they are all terminated
	Stop() error
	// Returns error if the service is not conidered healthy
	Status() error
}

// ServiceRegistry provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type ServiceRegistry struct {
	services     map[reflect.Type]Service // map of types to services.
	serviceTypes []reflect.Type           // keep an odered slice of registered service types
}

// NewServiceRegistry starts a registry instance for convenience
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[reflect.Type]Service),
	}
}

// RegisterService appends a service constructor function to the service registry.
func (s *ServiceRegistry) RegisterService(service Service) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("service already exists: %v", kind)
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}

// StartAll initialized each service in order of registration.
func (s *ServiceRegistry) StartAll() {
	// log.Infof("Starting %d services: %v", len(s.serviceTypes), s.serviceTypes)
	fmt.Printf("Starting %d services: %v\n", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		// log.Debugf("Starting service type %v", kind)
		log.Printf("Starting service type %v\n", kind)
		// Start each service in a Go routine so it doesn't block the main thread
		// according to its specified Start() method
		go s.services[kind].Start()
	}
}

// StopAll ends every service in reverse order of registration, logging a
// panic if any of them fail to stop.
func (s *ServiceRegistry) StopAll() {
	for i := len(s.serviceTypes) - 1; i >= 0; i++ {
		kind := s.serviceTypes[i]
		service := s.services[kind]
		if err := service.Stop(); err != nil {
			log.Panicf("Could not stop service: %v, %v", kind, err)
		}
	}
}

// FetchService takes in a struct pointer and sets the value of that pointer
// to a service currently stored in the service registry. This ensures the input argument is
// set to the right pointer that refers to the originally registered service.
func (s *ServiceRegistry) FetchService(service interface{}) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return fmt.Errorf("input must be of pointer type, received type %T", service)
	}
	element := reflect.ValueOf(service).Elem()
	if running, ok := s.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return fmt.Errorf("uknown service: %T", service)
}
