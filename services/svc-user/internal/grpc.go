package internal

import (
	"ops-monorepo/services/svc-user/config"
	"ops-monorepo/shared-libs/logger"
	userv1 "pb_schemas/user/v1"

	"google.golang.org/grpc"
)

type grpcServer struct {
	Server *grpc.Server
	user   *userImpl
	Log    logger.Logger
}

func NewGrpcServer(cfg *config.Config) *grpcServer {

	dep := InitDependencies(cfg)

	s := grpc.NewServer()

	return &grpcServer{
		Server: s,
		user:   &dep.Impl.userImpl,
		Log:    dep.log,
	}
}

func (s *grpcServer) Register() {

	// user implementation
	userv1.RegisterUserServiceServer(s.Server, s.user.handler)
}