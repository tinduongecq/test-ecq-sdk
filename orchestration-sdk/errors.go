package orchestrationsdk

import (
	"errors"
	"fmt"
)

// Validation errors
var (
	ErrNameRequired           = errors.New("name is required")
	ErrTemplateIDRequired     = errors.New("template_id is required")
	ErrResourcePoolIDRequired = errors.New("resource_pool_id is required")
	ErrInvalidCPU             = errors.New("cpu must be greater than 0")
	ErrInvalidMemory          = errors.New("memory is required")
	ErrInvalidDisk            = errors.New("disk is required")
)

// Client errors
var (
	ErrClientNotInitialized = errors.New("client not initialized")
	ErrInvalidBaseURL       = errors.New("invalid base URL")
	ErrRequestFailed        = errors.New("request failed")
	ErrResponseParseFailed  = errors.New("failed to parse response")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrNotFound             = errors.New("resource not found")
	ErrConflict             = errors.New("resource conflict")
	ErrInternalServerError  = errors.New("internal server error")
	ErrTimeout              = errors.New("request timeout")
)

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s - %s", e.StatusCode, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Code, e.Message)
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// NewAPIErrorWithDetails creates a new APIError with details
func NewAPIErrorWithDetails(statusCode int, code, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// RetryableError indicates an error that can be retried
type RetryableError struct {
	Err        error
	RetryAfter int // seconds
}

// Error implements the error interface
func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error (retry after %d seconds): %v", e.RetryAfter, e.Err)
}

// Unwrap returns the underlying error
func (e *RetryableError) Unwrap() error {
	return e.Err
}

// NewRetryableError creates a new RetryableError
func NewRetryableError(err error, retryAfter int) *RetryableError {
	return &RetryableError{
		Err:        err,
		RetryAfter: retryAfter,
	}
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var retryableErr *RetryableError
	return errors.As(err, &retryableErr)
}

// IsAPIError checks if an error is an API error
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// GetAPIError extracts APIError from an error
func GetAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	apiErr := GetAPIError(err)
	return apiErr != nil && apiErr.StatusCode == 404
}

// IsUnauthorized checks if an error is an unauthorized error
func IsUnauthorized(err error) bool {
	if errors.Is(err, ErrUnauthorized) {
		return true
	}
	apiErr := GetAPIError(err)
	return apiErr != nil && apiErr.StatusCode == 401
}

// IsConflict checks if an error is a conflict error
func IsConflict(err error) bool {
	if errors.Is(err, ErrConflict) {
		return true
	}
	apiErr := GetAPIError(err)
	return apiErr != nil && apiErr.StatusCode == 409
}
