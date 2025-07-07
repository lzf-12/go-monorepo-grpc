package handler

import (
	"context"
	"ops-monorepo/services/svc-user/internal/usecase"
	grpcErr "ops-monorepo/shared-libs/grpc/errors"
	"ops-monorepo/shared-libs/logger"
	userv1 "pb_schemas/user/v1"
)

type IUserHandler interface {
	ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error)
}

type userHandler struct {
	log            logger.Logger
	usecase        usecase.IUserUsecase
	grpcErrHandler *grpcErr.GRPCErrorHandler
	userv1.UnimplementedUserServiceServer
}

func NewUserHandler(log logger.Logger, usecase usecase.IUserUsecase, grpcErrHandler *grpcErr.GRPCErrorHandler) IUserHandler {
	return &userHandler{
		log:            log,
		usecase:        usecase,
		grpcErrHandler: grpcErrHandler,
	}
}

func (h *userHandler) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	h.log.Infof("received validate token request")

	if req.Token == "" {
		return &userv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Validate token
	user, err := h.usecase.ValidateToken(ctx, req.Token)
	if err != nil {
		h.log.Errorf("failed to validate token: %v", err)
		return &userv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &userv1.ValidateTokenResponse{
		Valid:     true,
		UserEmail: user.Email,
		Roles:     user.Roles,
	}, nil
}