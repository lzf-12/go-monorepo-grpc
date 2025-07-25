package errlib

const (
	ErrCodeEmailAlreadyUsed       string = "EMAIL_ALREADY_USED"
	ErrCodeUserNotFound           string = "USER_NOT_FOUND"
	ErrCodeInvalidEmailOrPassword string = "INVALID_EMAIL_OR_PASSWORD"
	ErrCodeTokenExpired           string = "TOKEN_EXPIRED"
	ErrCodeInvalidRefreshToken    string = "INVALID_REFRESH_TOKEN"
	ErrCodeInvalidInput           string = "INVALID_INPUT"
	ErrCodeAccessDenied           string = "ACCESS_DENIED"
	ErrCodeUnauthorized           string = "UNAUTHORIZED"
	ErrCodeForbidden              string = "FORBIDDEN"
	ErrCodeInternalServer         string = "INTERNAL_SERVER_ERROR"
	ErrCodeRateLimited            string = "RATE_LIMITED"
	ErrCodeValidation             string = "VALIDATION_ERROR"
	ErrCodeJSONUnmarshal          string = "JSON_UNMARSHAL_ERROR"
	ErrCodeJSONSyntax             string = "JSON_SYNTAX_ERROR"
	ErrCodeJSONBinding            string = "JSON_BINDING_ERROR"

	// db
	ErrCodeDBConnection    string = "DATABASE_CONNECTION_ERROR"
	ErrCodeDBQuery         string = "DATABASE_QUERY_ERROR"
	ErrCodeDBTransaction   string = "DATABASE_TRANSACTION_ERROR"
	ErrCodeDBConstraint    string = "DATABASE_CONSTRAINT_ERROR"
	ErrCodeDBDuplicate     string = "DATABASE_DUPLICATE_ERROR"
	ErrCodeDBTimeout       string = "DATABASE_TIMEOUT_ERROR"
	ErrCodeStorageNotFound string = "STORAGE_NOT_FOUND"
	ErrCodeStorageAccess   string = "STORAGE_ACCESS_ERROR"
	ErrCodeDataNotFound    string = "DATA_NOT_FOUND"

	// invetory
	ErrCodeReservationStock string = "FAILED_RESERVE_STOCK"
	ErrCodeReleaseStock     string = "FAILED_RELEASE_STOCK"
)
