package internal

import (
	"ops-monorepo/services/svc-inventory/config"
	"ops-monorepo/shared-libs/logger"
	inventoryv1 "pb_schemas/inventory/v1"

	"google.golang.org/grpc"
)

type grpcServer struct {
	Server    *grpc.Server
	inventory *inventoryImpl
	Log       logger.Logger
}

func NewGrpcServer(cfg *config.Config) *grpcServer {

	dep := InitDependencies(cfg)

	s := grpc.NewServer()

	return &grpcServer{
		Server:    s,
		inventory: &dep.Impl.inventoryImpl,
		Log:       dep.log,
	}
}

func (s *grpcServer) Register() {

	// inventory implementation
	inventoryv1.RegisterInventoryServiceServer(s.Server, s.inventory.handler)
}
