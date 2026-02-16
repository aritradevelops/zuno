# Error Handling Implementation

This document describes the comprehensive layer-by-layer error handling system implemented for the User API endpoints in the Goserve framework.

## Overview

The error handling system follows a clean architecture pattern with proper error propagation from the repository layer through the service layer to the handler layer. Each layer has its own error types and handling mechanisms while maintaining consistency in error responses.

## Error Handling Layers

### 1. Repository Layer (`internal/repository/errors.go`)

The repository layer defines low-level database and data access errors:

#### Error Types
- `RepositoryError`: Base error type for repository operations
- `RepositoryErrorCode`: Enum for different error categories

#### Error Categories
- **Database Operations**: `ErrDatabaseConnection`, `ErrDatabaseTimeout`, `ErrDatabaseQuery`, `ErrDatabaseInsert`, `ErrDatabaseUpdate`, `ErrDatabaseDelete`, `ErrDatabaseNotFound`
- **Data Validation**: `ErrInvalidData`, `ErrDuplicateEntry`, `ErrConstraintViolation`
- **Access Control**: `ErrAccessDenied`, `ErrUnauthorized`

#### Helper Functions
- `NewDatabaseQueryError()`: Database query failures
- `NewDatabaseInsertError()`: Database insert failures
- `NewDatabaseUpdateError()`: Database update failures
- `NewDatabaseDeleteError()`: Database delete failures
- `NewNotFoundError()`: Resource not found
- `NewDuplicateEntryError()`: Duplicate key violations
- `NewAccessDeniedError()`: Permission denied
- `NewInvalidDataError()`: Data validation failures

### 2. Service Layer (`internal/service/errors.go`)

The service layer defines business logic errors and provides error conversion from repository errors:

#### Error Interface
```go
type Error interface {
    Long() string
    Short() string
    Code() int
    error
}
```

#### Error Categories
- **Resource Operations**: `ErrResourceNotFound`, `ErrResourceAlreadyExists`, `ErrResourceInvalidData`, `ErrResourceAccessDenied`, `ErrResourceDeleteFailed`, `ErrResourceUpdateFailed`, `ErrResourceCreateFailed`, `ErrResourceListFailed`, `ErrResourceRestoreFailed`
- **General Service**: `ErrServiceInternal`, `ErrServiceTimeout`, `ErrServiceUnavailable`

#### Error Conversion
- `ConvertRepositoryError()`: Converts repository errors to service errors
- `AsServiceError()`: Type assertion helper for service errors

#### Key Design Principles
- **No HTTP Dependencies**: Service layer is transport-agnostic
- **Generic Resource Errors**: User errors are built on generic resource functions
- **Error Code Mapping**: HTTP status codes are determined by the transport layer

### 3. Handler Layer (`internal/transports/http/handler/user_handler.go`)

The handler layer acts as a pass-through for service errors:

#### Design Principles
- **Error Pass-Through**: Service errors are returned as-is to middleware
- **No HTTP Status Mapping**: Handlers don't determine HTTP status codes
- **Clean Separation**: Business logic errors remain separate from transport concerns

#### Error Response Format
```json
{
  "message": "Error description",
  "data": null,
  "info": null,
  "error": "Detailed error message"
}
```

### 4. Error Handler (`internal/transports/http/error_handler.go`)

The global error handler provides centralized error processing:

#### Error Processing Order
1. **Validation Errors**: Handle struct validation errors
2. **Service Errors**: Handle service layer errors with HTTP status mapping
3. **Repository Errors**: Handle repository errors with automatic HTTP status mapping
4. **Fiber Errors**: Handle framework errors
5. **Generic Errors**: Handle unexpected errors

#### HTTP Status Code Mapping Function
```go
func getHTTPStatusForServiceError(err service.Error) int {
    switch err.Code() {
    case service.ErrResourceNotFound:
        return httpb.StatusNotFound
    case service.ErrResourceAlreadyExists:
        return httpb.StatusConflict
    case service.ErrResourceInvalidData:
        return httpb.StatusUnprocessableEntity
    case service.ErrResourceAccessDenied:
        return httpb.StatusForbidden
    case service.ErrResourceDeleteFailed, service.ErrResourceUpdateFailed, service.ErrResourceCreateFailed, service.ErrResourceListFailed, service.ErrResourceRestoreFailed:
        return httpb.StatusInternalServerError
    case service.ErrServiceInternal, service.ErrServiceTimeout, service.ErrServiceUnavailable:
        return httpb.StatusInternalServerError
    default:
        return httpb.StatusInternalServerError
    }
}
```

## Error Flow Example

### Create User Operation

1. **Repository Layer**:
   ```go
   // MongoDB adapter
   if _, err := r.collection.InsertOne(ctx, user); err != nil {
       return nil, repository.NewDatabaseInsertError("create user", err)
   }
   ```

2. **Service Layer**:
   ```go
   // Service layer
   user, err := s.userRepository.Create(ctx, actor, repository.UserFields(payload))
   if err != nil {
       return nil, service.ConvertRepositoryError(err)
   }
   ```

3. **Handler Layer**:
   ```go
   // Handler layer
   user, err := h.userService.Create(c.Context(), actor, service.UserFields(payload))
   if err != nil {
       var serviceErr service.Error
       if service.AsServiceError(err, &serviceErr) {
           return err // Return service error as-is for middleware handling
       }
   }
   ```

4. **Error Handler**:
   ```go
   // Global error handler
   var serviceErr service.Error
   if service.AsServiceError(e, &serviceErr) {
       httpStatus := getHTTPStatusForServiceError(serviceErr)
       c.Status(httpStatus)
       return c.JSON(handler.Failure(translation.Localize(c, fmt.Sprintf("errors.%d", httpStatus)), serviceErr.Long()))
   }
   ```

## Error Response Examples

### 404 Not Found
```json
{
  "message": "User not found",
  "data": null,
  "info": null,
  "error": "User with ID 12345 not found"
}
```

### 409 Conflict
```json
{
  "message": "User already exists",
  "data": null,
  "info": null,
  "error": "Duplicate entry for User in field 'email' with value 'user@example.com'"
}
```

### 422 Validation Error
```json
{
  "message": "Invalid user data",
  "data": null,
  "info": null,
  "error": "The email is required"
}
```

### 500 Internal Server Error
```json
{
  "message": "Internal server error",
  "data": null,
  "info": null,
  "error": "Service Error [Code: 2001, HTTP: 500]: Internal service error"
}
```

## Error Localization

Error messages are localized using the i18n system:

### English (en.yaml)
```yaml
errors:
  "400": "Invalid request format or parameters"
  "401": "Authentication required"
  "403": "Access denied"
  "404": "Resource not found"
  "409": "Conflict - resource already exists"
  "422": "Validation failed"
  "500": "Internal server error"
```

### Bengali (bn.yaml)
```yaml
errors:
  "400": "অবৈধ অনুরোধ বিন্যাস বা প্যারামিটার"
  "401": "প্রমাণীকরণ প্রয়োজন"
  "403": "অ্যাক্সেস অস্বীকৃত"
  "404": "সম্পদ পাওয়া যায়নি"
  "409": "দ্বন্দ্ব - সম্পদ ইতিমধ্যে অস্তিত্ব"
  "422": "বৈধকরণ ব্যর্থ"
  "500": "অভ্যন্তরীণ সার্ভার ত্রুটি"
```

## Best Practices

### 1. Error Wrapping
Always wrap lower-level errors with context:
```go
return nil, repository.NewDatabaseInsertError("create user", err)
```

### 2. Error Conversion
Convert errors at layer boundaries:
```go
return nil, service.ConvertRepositoryError(err)
```

### 3. HTTP Status Codes
HTTP status codes are determined by the transport layer based on service error codes:
- 400 for bad requests (handled by framework)
- 401 for authentication failures (handled by middleware)
- 403 for authorization failures (mapped from `ErrResourceAccessDenied`)
- 404 for not found (mapped from `ErrResourceNotFound`)
- 409 for conflicts (mapped from `ErrResourceAlreadyExists`)
- 422 for validation errors (mapped from `ErrResourceInvalidData`)
- 500 for internal errors (mapped from service operation failures)

### 4. Error Messages
Provide clear, actionable error messages:
- Include resource names and IDs
- Specify the operation that failed
- Provide context for debugging

### 5. Logging
Log errors with appropriate levels:
```go
logger.Error().Err(err).Msg("failed to create user")
```

## Testing Error Scenarios

### Unit Tests
Test error handling in each layer:
```go
func TestUserService_Create_Error(t *testing.T) {
    // Test repository error conversion
    // Test service error handling
    // Test handler error responses
}
```

### Integration Tests
Test end-to-end error scenarios:
```go
func TestUserAPI_Create_Error(t *testing.T) {
    // Test HTTP error responses
    // Test error message localization
    // Test error status codes
}
```

## Monitoring and Observability

### Error Metrics
Track error rates and types:
- Error count by HTTP status code
- Error count by operation
- Error count by error type

### Error Logging
Structured logging for debugging:
```json
{
  "level": "error",
  "message": "failed to create user",
  "error": "duplicate key error",
  "operation": "create user",
  "user_id": "12345",
  "timestamp": "2026-02-15T17:56:00Z"
}
```

This comprehensive error handling system ensures consistent error responses, proper HTTP status codes, and maintainable error handling across all layers of the application.