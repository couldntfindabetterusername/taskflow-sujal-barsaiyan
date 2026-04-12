package errors

import (
 "errors"
 "fmt"
)

var (
 // ErrNotFound is returned when a resource is not found
 ErrNotFound = errors.New("resource not found")

 // ErrDuplicate is returned when a resource already exists
 ErrDuplicate = errors.New("resource already exists")

 // ErrInvalidCredentials is returned when authentication credentials are invalid
 ErrInvalidCredentials = errors.New("invalid credentials")

 // ErrValidation is returned when input validation fails
 ErrValidation = errors.New("validation error")

 // ErrUnauthorized is returned when a user is not authorized to perform an action
 ErrUnauthorized = errors.New("unauthorized")

 // ErrInvalidToken is returned when a JWT token is invalid
 ErrInvalidToken = errors.New("invalid token")
)

// ValidationError represents a validation error with details
type ValidationError struct {
 Field   string
 Message string
}

func (e *ValidationError) Error() string {
 return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
 return &ValidationError{
 Field:   field,
 Message: message,
 }
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
 return errors.Is(err, ErrNotFound)
}

// IsDuplicate checks if an error is a duplicate error
func IsDuplicate(err error) bool {
 return errors.Is(err, ErrDuplicate)
}

// IsInvalidCredentials checks if an error is an invalid credentials error
func IsInvalidCredentials(err error) bool {
 return errors.Is(err, ErrInvalidCredentials)
}

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
 if errors.Is(err, ErrValidation) {
 return true
 }
 var validationErr *ValidationError
 return errors.As(err, &validationErr)
}

// IsUnauthorized checks if an error is an unauthorized error
func IsUnauthorized(err error) bool {
 return errors.Is(err, ErrUnauthorized)
}

// IsInvalidToken checks if an error is an invalid token error
func IsInvalidToken(err error) bool {
 return errors.Is(err, ErrInvalidToken)
}
