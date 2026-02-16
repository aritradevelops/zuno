package repository

import (
	"errors"
	"fmt"
	"time"
)

// RepositoryError represents errors that occur at the repository/database layer
type RepositoryError struct {
	Code    RepositoryErrorCode
	Message string
	Cause   error
	Details map[string]any
}

type RepositoryErrorCode string

const (
	// Database connection and operation errors
	ErrDatabaseConnection RepositoryErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrDatabaseTimeout    RepositoryErrorCode = "DATABASE_TIMEOUT_ERROR"
	ErrDatabaseQuery      RepositoryErrorCode = "DATABASE_QUERY_ERROR"
	ErrDatabaseInsert     RepositoryErrorCode = "DATABASE_INSERT_ERROR"
	ErrDatabaseUpdate     RepositoryErrorCode = "DATABASE_UPDATE_ERROR"
	ErrDatabaseDelete     RepositoryErrorCode = "DATABASE_DELETE_ERROR"
	ErrDatabaseNotFound   RepositoryErrorCode = "DATABASE_NOT_FOUND_ERROR"

	// Data validation errors
	ErrInvalidData         RepositoryErrorCode = "INVALID_DATA_ERROR"
	ErrDuplicateEntry      RepositoryErrorCode = "DUPLICATE_ENTRY_ERROR"
	ErrConstraintViolation RepositoryErrorCode = "CONSTRAINT_VIOLATION_ERROR"

	// Permission and access errors
	ErrAccessDenied RepositoryErrorCode = "ACCESS_DENIED_ERROR"
	ErrUnauthorized RepositoryErrorCode = "UNAUTHORIZED_ERROR"
)

func (e *RepositoryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Repository Error [%s]: %s. Caused by: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("Repository Error [%s]: %s", e.Code, e.Message)
}

func (e *RepositoryError) Unwrap() error {
	return e.Cause
}

// Helper functions to create specific repository errors
func NewDatabaseConnectionError(cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseConnection,
		Message: "Failed to connect to database",
		Cause:   cause,
		Details: map[string]any{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewDatabaseTimeoutError(operation string, timeout time.Duration, cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseTimeout,
		Message: fmt.Sprintf("Database operation '%s' timed out after %v", operation, timeout),
		Cause:   cause,
		Details: map[string]any{
			"operation": operation,
			"timeout":   timeout.String(),
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewNotFoundError(resource, id string) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseNotFound,
		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
		Details: map[string]any{
			"resource":  resource,
			"id":        id,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewDuplicateEntryError(resource, field string, value any) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDuplicateEntry,
		Message: fmt.Sprintf("Duplicate entry for %s in field '%s' with value '%v'", resource, field, value),
		Details: map[string]any{
			"resource":  resource,
			"field":     field,
			"value":     value,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewAccessDeniedError(resource, action string, actorID string) *RepositoryError {
	return &RepositoryError{
		Code:    ErrAccessDenied,
		Message: fmt.Sprintf("Access denied: %s %s for actor %s", action, resource, actorID),
		Details: map[string]any{
			"resource":  resource,
			"action":    action,
			"actor_id":  actorID,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewInvalidDataError(message string, details map[string]any) *RepositoryError {
	return &RepositoryError{
		Code:    ErrInvalidData,
		Message: message,
		Details: details,
	}
}

// IsRepositoryError checks if an error is a RepositoryError
func IsRepositoryError(err error) bool {
	_, ok := err.(*RepositoryError)
	return ok
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var repoErr *RepositoryError
	if !errors.As(err, &repoErr) {
		return false
	}
	return repoErr.Code == ErrDatabaseNotFound
}

// IsDuplicateEntryError checks if an error is a duplicate entry error
func IsDuplicateEntryError(err error) bool {
	var repoErr *RepositoryError
	if !errors.As(err, &repoErr) {
		return false
	}
	return repoErr.Code == ErrDuplicateEntry
}

// IsAccessDeniedError checks if an error is an access denied error
func IsAccessDeniedError(err error) bool {
	var repoErr *RepositoryError
	if !errors.As(err, &repoErr) {
		return false
	}
	return repoErr.Code == ErrAccessDenied
}

// Helper functions for specific database operations
func NewDatabaseQueryError(operation string, cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseQuery,
		Message: fmt.Sprintf("Database query failed for operation '%s'", operation),
		Cause:   cause,
		Details: map[string]any{
			"operation": operation,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewDatabaseInsertError(operation string, cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseInsert,
		Message: fmt.Sprintf("Database insert failed for operation '%s'", operation),
		Cause:   cause,
		Details: map[string]any{
			"operation": operation,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewDatabaseUpdateError(operation string, cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseUpdate,
		Message: fmt.Sprintf("Database update failed for operation '%s'", operation),
		Cause:   cause,
		Details: map[string]any{
			"operation": operation,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}

func NewDatabaseDeleteError(operation string, cause error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrDatabaseDelete,
		Message: fmt.Sprintf("Database delete failed for operation '%s'", operation),
		Cause:   cause,
		Details: map[string]any{
			"operation": operation,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
}
