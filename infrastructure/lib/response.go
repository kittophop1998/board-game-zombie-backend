package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information in response
type ErrorInfo struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// ResponseSuccess sends a successful response with data
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// ResponseSuccessWithMessage sends a successful response with data and message
func ResponseSuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ResponseCreated sends a created response (201) with data
func ResponseCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// ResponseCreatedWithMessage sends a created response with data and message
func ResponseCreatedWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ResponseNoContent sends a no content response (204)
func ResponseNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ResponseError sends an error response with specified status code
func ResponseError(c *gin.Context, statusCode int, errorType ErrorType, message string, details ...string) {
	errorInfo := &ErrorInfo{
		Type:    errorType,
		Message: message,
	}

	if len(details) > 0 {
		errorInfo.Details = details[0]
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorInfo,
	})
}

// ResponseBadRequest sends a bad request error (400)
func ResponseBadRequest(c *gin.Context, message string, details ...string) {
	ResponseError(c, http.StatusBadRequest, ErrorTypeValidation, message, details...)
}

// ResponseUnauthorized sends an unauthorized error (401)
func ResponseUnauthorized(c *gin.Context, message string, details ...string) {
	if message == "" {
		message = "Authentication required"
	}
	ResponseError(c, http.StatusUnauthorized, ErrorTypeAuthentication, message, details...)
}

// ResponseForbidden sends a forbidden error (403)
func ResponseForbidden(c *gin.Context, message string, details ...string) {
	if message == "" {
		message = "Access forbidden"
	}
	ResponseError(c, http.StatusForbidden, ErrorTypeAuthorization, message, details...)
}

// ResponseNotFound sends a not found error (404)
func ResponseNotFound(c *gin.Context, message string, details ...string) {
	if message == "" {
		message = "Resource not found"
	}
	ResponseError(c, http.StatusNotFound, ErrorTypeNotFound, message, details...)
}

// ResponseConflict sends a conflict error (409)
func ResponseConflict(c *gin.Context, message string, details ...string) {
	ResponseError(c, http.StatusConflict, ErrorTypeConflict, message, details...)
}

// ResponseInternalServerError sends an internal server error (500)
func ResponseInternalServerError(c *gin.Context, message string, details ...string) {
	if message == "" {
		message = "Internal server error"
	}
	ResponseError(c, http.StatusInternalServerError, ErrorTypeInternal, message, details...)
}

// ResponseValidationError sends a validation error (400) with validation details
func ResponseValidationError(c *gin.Context, validationErrors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Type:    ErrorTypeValidation,
			Message: "Validation failed",
		},
		Data: validationErrors,
	})
}

// ResponseCustom sends a custom response with full control
func ResponseCustom(c *gin.Context, statusCode int, response Response) {
	c.JSON(statusCode, response)
}
