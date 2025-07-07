package grpc_errors

// BusinessError represents a business logic error with context
type AppError struct {
	Type    ErrorType
	Message string
	Details map[string]interface{}
}

func (e AppError) Error() string {
	return e.Message
}

// NewBusinessError creates a new business error
func NewAppError(errType ErrorType, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Details: details,
	}
}
