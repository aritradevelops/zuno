package service

import (
	"errors"
	"fmt"
	"goserve/internal/repository"
	"time"
)

// ServiceError represents errors that occur at the business logic layer
type ServiceError struct {
	Code    ServiceErrorCode
	Message string
	Cause   error
	Details map[string]any
}

type ServiceErrorCode string

const (
	// Validation
	ErrValidationFailed ServiceErrorCode = "VALIDATION_FAILED"

	// State & lifecycle errors
	ErrAlreadyExists    ServiceErrorCode = "ALREADY_EXISTS"
	ErrNotFound         ServiceErrorCode = "NOT_FOUND"
	ErrOperationBlocked ServiceErrorCode = "OPERATION_BLOCKED"

	// Authorization & permissions
	ErrUnauthorized ServiceErrorCode = "UNAUTHORIZED"
	ErrForbidden    ServiceErrorCode = "FORBIDDEN"

	// Dependencies & orchestration
	ErrExternalService ServiceErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrTimeout         ServiceErrorCode = "SERVICE_TIMEOUT"

	// Internal
	ErrInternal ServiceErrorCode = "INTERNAL_SERVICE_ERROR"
)

func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Service Error [%s]: %s. Caused by: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("Service Error [%s]: %s", e.Code, e.Message)
}

func (e *ServiceError) Unwrap() error {
	return e.Cause
}
func NewValidationError(message string, details map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrValidationFailed,
		Message: message,
		Details: details,
	}
}

func NewAlreadyExistsError(resource, key string) *ServiceError {
	return &ServiceError{
		Code:    ErrAlreadyExists,
		Message: fmt.Sprintf("%s already exists with key %s", resource, key),
		Details: map[string]any{
			"resource": resource,
			"key":      key,
		},
	}
}

func NewServiceNotFoundError(resource, id string) *ServiceError {
	return &ServiceError{
		Code:    ErrNotFound,
		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
		Details: map[string]any{
			"resource": resource,
			"id":       id,
		},
	}
}

func NewOperationBlockedError(reason string, details map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrOperationBlocked,
		Message: reason,
		Details: details,
	}
}
func NewUnauthorizedError() *ServiceError {
	return &ServiceError{
		Code:    ErrUnauthorized,
		Message: "Unauthorized",
	}
}

func NewForbiddenError(actor, action, resource string) *ServiceError {
	return &ServiceError{
		Code:    ErrForbidden,
		Message: fmt.Sprintf("Actor %s cannot %s %s", actor, action, resource),
		Details: map[string]any{
			"actor":    actor,
			"action":   action,
			"resource": resource,
		},
	}
}

func NewExternalServiceError(service string, cause error) *ServiceError {
	return &ServiceError{
		Code:    ErrExternalService,
		Message: fmt.Sprintf("External service error: %s", service),
		Cause:   cause,
	}
}

func NewServiceTimeout(operation string, timeout time.Duration) *ServiceError {
	return &ServiceError{
		Code:    ErrTimeout,
		Message: fmt.Sprintf("Service operation '%s' timed out after %v", operation, timeout),
		Details: map[string]any{
			"operation": operation,
			"timeout":   timeout.String(),
		},
	}
}
func NewInternalError(message string, cause error) *ServiceError {
	return &ServiceError{
		Code:    ErrInternal,
		Message: message,
		Cause:   cause,
	}
}
func IsServiceError(err error) bool {
	var se *ServiceError
	return errors.As(err, &se)
}

func HasCode(err error, code ServiceErrorCode) bool {
	var se *ServiceError
	if !errors.As(err, &se) {
		return false
	}
	return se.Code == code
}

func IsNotFound(err error) bool {
	return HasCode(err, ErrNotFound)
}

func IsValidationError(err error) bool {
	return HasCode(err, ErrValidationFailed)
}

func IsForbidden(err error) bool {
	return HasCode(err, ErrForbidden)
}

func ConvertRepositoryError(err error) error {
	var repoErr *repository.RepositoryError
	if errors.As(err, &repoErr) {
		switch repoErr.Code {
		case repository.ErrAccessDenied:
			return NewUnauthorizedError()
		}
	}
	return NewInternalError("unknown error occurred", err)
}
