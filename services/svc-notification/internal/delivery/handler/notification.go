package handler

import (
	"context"
	"ops-monorepo/services/svc-notification/internal/usecase"
	grpcErr "ops-monorepo/shared-libs/grpc/errors"
	"ops-monorepo/shared-libs/logger"
	notificationv1 "pb_schemas/notification/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type INotificationHandler interface {
	SendEmail(ctx context.Context, req *notificationv1.SendEmailRequest) (*notificationv1.SendEmailResponse, error)
}

type notificationHandler struct {
	log            logger.Logger
	usecase        usecase.INotificationUsecase
	grpcErrHandler *grpcErr.GRPCErrorHandler
	notificationv1.UnimplementedNotificationServiceServer
}

func NewNotificationHandler(log logger.Logger, usecase usecase.INotificationUsecase, grpcErrHandler *grpcErr.GRPCErrorHandler) INotificationHandler {
	return &notificationHandler{
		log:            log,
		usecase:        usecase,
		grpcErrHandler: grpcErrHandler,
	}
}

func (h *notificationHandler) SendEmail(ctx context.Context, req *notificationv1.SendEmailRequest) (*notificationv1.SendEmailResponse, error) {
	h.log.Infof("received send email request for: %s", req.To)

	// Validate request
	if req.To == "" {
		return &notificationv1.SendEmailResponse{
			Success: false,
			Message: "recipient email is required",
		}, nil
	}

	if req.Subject == "" {
		return &notificationv1.SendEmailResponse{
			Success: false,
			Message: "subject is required",
		}, nil
	}

	if req.Body == "" {
		return &notificationv1.SendEmailResponse{
			Success: false,
			Message: "body is required",
		}, nil
	}

	// Send email
	emailLog, err := h.usecase.SendEmail(ctx, req.To, req.Subject, req.Body, req.From, req.IsHtml)
	if err != nil {
		h.log.Errorf("failed to send email: %v", err)
		return &notificationv1.SendEmailResponse{
			Success:   false,
			Message:   err.Error(),
			EmailId:   emailLog.ID,
			Timestamp: timestamppb.New(emailLog.UpdatedAt),
		}, nil
	}

	return &notificationv1.SendEmailResponse{
		Success:   true,
		Message:   "email sent successfully",
		EmailId:   emailLog.ID,
		Timestamp: timestamppb.New(emailLog.UpdatedAt),
	}, nil
}