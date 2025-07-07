package grpc

import (
	"log"

	// pb_schemas generated interface
	inventoryv1 "pb_schemas/inventory/v1"
	userv1 "pb_schemas/user/v1"
)

// aliases
type (
	InvClient  = inventoryv1.InventoryServiceClient
	UserClient = userv1.UserServiceClient
)

type ServiceClients struct {
	registry *ClientRegistry
}

func NewServiceClients() *ServiceClients {
	return &ServiceClients{
		registry: NewClientRegistry(),
	}
}

func (s *ServiceClients) Close() {
	s.registry.Close()
}

// inventory service

func (r *ClientRegistry) GetInventoryClient(target string) (inventoryv1.InventoryServiceClient, error) {
	conn, err := r.GetConnection(target)
	if err != nil {
		return nil, err
	}
	return inventoryv1.NewInventoryServiceClient(conn), nil
}

func (s *ServiceClients) Inventory(target string) *inventoryv1.InventoryServiceClient {
	client, err := s.registry.GetInventoryClient(target)
	if err != nil {
		log.Fatalf("Failed to get inventory client for %s: %v", target, err)
	}
	return &client
}

// user service

type UserGrpcClientInterface = userv1.UserServiceClient

func (r *ClientRegistry) GetUserClient(target string) (userv1.UserServiceClient, error) {
	conn, err := r.GetConnection(target)
	if err != nil {
		return nil, err
	}
	return userv1.NewUserServiceClient(conn), nil
}

func (s *ServiceClients) User(target string) userv1.UserServiceClient {
	client, err := s.registry.GetUserClient(target)
	if err != nil {
		log.Fatalf("Failed to get user client for %s: %v", target, err)
	}
	return client
}

// add other service here
// make sure the grpc method already generated on ./pb_schemas directory
