package internal

import (
	"ops-monorepo/services/svc-notification/config"
	"ops-monorepo/shared-libs/logger"
	notificationv1 "pb_schemas/notification/v1"

	"google.golang.org/grpc"
)

type grpcServer struct {
	Server       *grpc.Server
	notification *notificationImpl
	Log          logger.Logger
}

func NewGrpcServer(cfg *config.Config) *grpcServer {

	dep := InitDependencies(cfg)

	s := grpc.NewServer()

	return &grpcServer{
		Server:       s,
		notification: &dep.Impl.notificationImpl,
		Log:          dep.log,
	}
}

func (s *grpcServer) Register() {

	// notification implementation
	notificationv1.RegisterNotificationServiceServer(s.Server, s.notification.handler)
}