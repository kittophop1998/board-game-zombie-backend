package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeConflict       ErrorType = "conflict"
	ErrorTypeInternal       ErrorType = "internal"
	ErrorTypeBusiness       ErrorType = "business"
)

// ProblemDetail represents the RFC 7807 Problem Details structure
type ProblemDetail struct {
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Status   int                    `json:"status"`
	Detail   string                 `json:"detail"`
	Instance string                 `json:"instance"`
	Errors   map[string]interface{} `json:"errors,omitempty"`
}

// AppError represents application-specific error
type AppError struct {
	Type       ErrorType
	StatusCode int
	Message    string
	Details    string
	Errors     map[string]interface{}
	Cause      error
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error
func NewAppError(errorType ErrorType, statusCode int, message string, details string) *AppError {
	return &AppError{
		Type:       errorType,
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
		Errors:     make(map[string]interface{}),
	}
}

// WithErrors adds validation errors to the AppError
func (e *AppError) WithErrors(errors map[string]interface{}) *AppError {
	e.Errors = errors
	return e
}

// WithCause adds underlying error cause
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// Common error constructors
func NewValidationError(detail string) *AppError {
	return NewAppError(ErrorTypeValidation, http.StatusBadRequest, "Validation Error", detail)
}

func NewAuthenticationError(detail string) *AppError {
	return NewAppError(ErrorTypeAuthentication, http.StatusUnauthorized, "Authentication Error", detail)
}

func NewAuthorizationError(detail string) *AppError {
	return NewAppError(ErrorTypeAuthorization, http.StatusForbidden, "Authorization Error", detail)
}

func NewNotFoundError(detail string) *AppError {
	return NewAppError(ErrorTypeNotFound, http.StatusNotFound, "Resource Not Found", detail)
}

func NewConflictError(detail string) *AppError {
	return NewAppError(ErrorTypeConflict, http.StatusConflict, "Conflict Error", detail)
}

func NewInternalError(detail string) *AppError {
	return NewAppError(ErrorTypeInternal, http.StatusInternalServerError, "Internal Server Error", detail)
}

func NewBusinessError(detail string) *AppError {
	return NewAppError(ErrorTypeBusiness, http.StatusBadRequest, "Business Logic Error", detail)
}

// ErrorHandler handles errors and returns appropriate HTTP response
func ErrorHandler(c *gin.Context, err error) {
	var appErr *AppError
	var statusCode int
	var problemDetail ProblemDetail

	// Check if it's an AppError
	if appError, ok := err.(*AppError); ok {
		appErr = appError
		statusCode = appErr.StatusCode
	} else {
		// Handle unknown errors as internal server error
		appErr = NewInternalError("An unexpected error occurred")
		statusCode = http.StatusInternalServerError
	}

	// Build the base URL for error types
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s/errors", scheme, c.Request.Host)

	// Create problem detail response
	problemDetail = ProblemDetail{
		Type:     fmt.Sprintf("%s/%s", baseURL, appErr.Type),
		Title:    appErr.Message,
		Status:   statusCode,
		Detail:   appErr.Details,
		Instance: c.Request.RequestURI,
	}

	// Add validation errors if they exist
	if len(appErr.Errors) > 0 {
		problemDetail.Errors = appErr.Errors
	}

	// Set content type
	c.Header("Content-Type", "application/problem+json")

	// Return JSON response
	c.JSON(statusCode, problemDetail)
}

// ErrorMiddleware is a Gin middleware for handling panics and errors
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Handle panic
				var appErr *AppError

				switch v := err.(type) {
				case *AppError:
					appErr = v
				case error:
					appErr = NewInternalError(v.Error())
				case string:
					appErr = NewInternalError(v)
				default:
					appErr = NewInternalError("An unexpected error occurred")
				}

				ErrorHandler(c, appErr)
				c.Abort()
				return
			}
		}()

		c.Next()

		// Handle errors set by handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			ErrorHandler(c, err.Err)
		}
	}
}

// ValidationErrorsFromJSON parses validation errors from JSON tags
func ValidationErrorsFromJSON(jsonErr error) map[string]interface{} {
	errors := make(map[string]interface{})

	if jsonErr == nil {
		return errors
	}

	switch v := jsonErr.(type) {
	case *json.UnmarshalTypeError:
		errors[v.Field] = []string{fmt.Sprintf("Invalid type. Expected %s", v.Type)}
	case *json.SyntaxError:
		errors["json"] = []string{"Invalid JSON syntax"}
	default:
		errors["json"] = []string{jsonErr.Error()}
	}

	return errors
}

// Helper function to create validation error with field errors
func NewValidationErrorWithFields(detail string, fieldErrors map[string]interface{}) *AppError {
	return NewValidationError(detail).WithErrors(fieldErrors)
}

// Example usage functions for common scenarios

// HandleBindingError handles JSON binding errors from Gin
func HandleBindingError(c *gin.Context, err error) {
	fieldErrors := ValidationErrorsFromJSON(err)
	appErr := NewValidationErrorWithFields("Request validation failed", fieldErrors)
	ErrorHandler(c, appErr)
}

// HandleDatabaseError handles database-related errors
func HandleDatabaseError(c *gin.Context, err error, operation string) {
	var appErr *AppError

	// You can add specific database error handling here
	// For example, check for specific PostgreSQL error codes

	if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
		appErr = NewNotFoundError(fmt.Sprintf("Resource not found during %s", operation))
	} else {
		appErr = NewInternalError(fmt.Sprintf("Database error during %s", operation)).WithCause(err)
	}

	ErrorHandler(c, appErr)
}
