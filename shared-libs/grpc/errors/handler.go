package grpc_errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCErrorHandler converts business errors to gRPC status errors
type GRPCErrorHandler struct {
	// You can add configuration here if needed
}

// NewGRPCErrorHandler creates a new error handler instance
func NewGRPCErrorHandler() *GRPCErrorHandler {
	return &GRPCErrorHandler{}
}

// HandleError converts various error types to appropriate gRPC status errors
func (h *GRPCErrorHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}

	// handle known/mapped AppError types
	if bizErr, ok := err.(*AppError); ok {
		return h.handleBusinessError(bizErr)
	}

	// Handle other error types or fallback to internal error
	return h.createStatusError(codes.Internal, "internal server error", nil)
}

// handleBusinessError converts BusinessError to gRPC status error
func (h *GRPCErrorHandler) handleBusinessError(bizErr *AppError) error {
	switch bizErr.Type {
	case ValidationError:
		return h.createStatusError(codes.InvalidArgument, bizErr.Message, bizErr.Details)

	// case SKUNotFound:
	// 	return h.createStatusError(codes.NotFound, bizErr.Message, bizErr.Details)

	// case SKUUOMPairMismatch:
	// 	return h.createStatusError(codes.InvalidArgument, bizErr.Message, bizErr.Details)

	case InsufficientQuantity:
		return h.createStatusError(codes.FailedPrecondition, bizErr.Message, bizErr.Details)

	case InsufficientReservedQuantity:
		return h.createStatusError(codes.FailedPrecondition, bizErr.Message, bizErr.Details)

	case DbError:
		return h.createStatusError(codes.Unavailable, "database operation failed", nil)

	case InternalServerError:
		return h.createStatusError(codes.Internal, "internal server error", nil)

	default:
		return h.createStatusError(codes.Internal, "unknown error", nil)
	}
}

// createStatusError creates a gRPC status error with optional details
func (h *GRPCErrorHandler) createStatusError(code codes.Code, message string, details map[string]interface{}) error {
	st := status.New(code, message)

	if details != nil && len(details) > 0 {
		// Convert details to Any proto message
		// You can customize this based on your proto definitions
		detailsProto := h.convertDetailsToProto(details)
		if detailsProto != nil {
			st, _ = st.WithDetails(detailsProto)
		}
	}

	return st.Err()
}

// TODO: implement details to proto
func (h *GRPCErrorHandler) convertDetailsToProto(details map[string]interface{}) *anypb.Any {
	return nil
}

func NewValidationError(message string, fieldErrors map[string]string) *AppError {
	details := make(map[string]interface{})
	if fieldErrors != nil {
		details["field_errors"] = fieldErrors
	}
	return NewAppError(ValidationError, message, details)
}

// func NewSKUNotFoundError(sku string) *AppError {
// 	return NewBusinessError(SKUNotFound, fmt.Sprintf("SKU '%s' not found", sku), map[string]interface{}{
// 		"sku": sku,
// 	})
// }

// func NewSKUUOMPairMismatchError(sku, uom string) *AppError {
// 	return NewBusinessError(SKUUOMPairMismatch, fmt.Sprintf("SKU '%s' does not support UOM '%s'", sku, uom), map[string]interface{}{
// 		"sku": sku,
// 		"uom": uom,
// 	})
// }

func NewInsufficientQuantityError(sku string, requested, available int64) *AppError {
	return NewAppError(InsufficientQuantity, fmt.Sprintf("insufficient quantity for SKU '%s': requested %d, available %d", sku, requested, available), map[string]interface{}{
		"sku":       sku,
		"requested": requested,
		"available": available,
	})
}

func NewInsufficientReservedQuantityError(sku string, requested, reserved int64) *AppError {
	return NewAppError(InsufficientReservedQuantity, fmt.Sprintf("insufficient reserved quantity for SKU '%s': requested %d, reserved %d", sku, requested, reserved), map[string]interface{}{
		"sku":       sku,
		"requested": requested,
		"reserved":  reserved,
	})
}

func NewDbError(operation string, err error) *AppError {
	return NewAppError(DbError, fmt.Sprintf("database operation '%s' failed", operation), map[string]interface{}{
		"operation": operation,
		"error":     err.Error(),
	})
}

func NewDbTransactionError(operation string, err error) *AppError {
	return NewAppError(DbTransactionError, fmt.Sprintf("db transaction operation '%s' failed", operation), map[string]interface{}{
		"operation": operation,
		"error":     err.Error(),
	})
}

func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "internal server error"
	}
	return NewAppError(InternalServerError, message, nil)
}
