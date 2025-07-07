package grpc_errors

// ErrorType represents different business error types
type ErrorType int

const (
	ValidationError ErrorType = iota
	InsufficientQuantity
	InsufficientReservedQuantity
	DbError
	DbTransactionError
	InternalServerError
)
