// internal/errors/errors.go
package errors

import "fmt"

type ErrorCode string

const (
	// Client Errors (4xx)
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"

	// Server Errors (5xx)
	ErrInternalServer    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDatabaseOperation ErrorCode = "DATABASE_ERROR"
	ErrExternalService   ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// Business Logic Errors
	ErrOrganizationInactive ErrorCode = "ORGANIZATION_INACTIVE"
	ErrInvalidStatus        ErrorCode = "INVALID_STATUS"
	ErrLimitExceeded        ErrorCode = "LIMIT_EXCEEDED"
)

type AppError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	HTTPStatus int         `json:"-"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error constructors
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrBadRequest,
		Message:    message,
		HTTPStatus: 400,
	}
}

func NewValidationError(details interface{}) *AppError {
	return &AppError{
		Code:       ErrValidation,
		Message:    "Validation failed",
		Details:    details,
		HTTPStatus: 400,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       ErrNotFound,
		Message:    message,
		HTTPStatus: 404,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:       ErrInternalServer,
		Message:    message,
		HTTPStatus: 500,
	}
}

// Add more error constructors as needed
