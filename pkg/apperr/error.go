package apperr

import (
	"fmt"
	"net/http"
)

// AppError defines a custom application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Cause      error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap provides compatibility with errors.Is/As
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError
func New(code string, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

// WithCause adds the underlying error
func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

// Error Codes
const (
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeInvalidInput   = "INVALID_INPUT"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeInternalError  = "INTERNAL_ERROR"
	ErrCodeUserNotFound   = "USER_NOT_FOUND"
	ErrCodeUserConflict   = "USER_ALREADY_EXISTS"
	ErrCodeFileTooLarge   = "FILE_TOO_LARGE"
	ErrCodeUploadFailed   = "UPLOAD_FAILED"
)

// Helper functions for common errors

func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message, http.StatusBadRequest)
}

func InvalidInput(message string) *AppError {
	return New(ErrCodeInvalidInput, message, http.StatusBadRequest)
}

func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

func NotFound(message string) *AppError {
	return New(ErrCodeNotFound, message, http.StatusNotFound)
}

func InternalServerError(err error) *AppError {
	return New(ErrCodeInternalError, "Internal Server Error", http.StatusInternalServerError).WithCause(err)
}

// Specific Domain Errors

func UserNotFound() *AppError {
	return New(ErrCodeUserNotFound, "User not found", http.StatusNotFound)
}

func UserAlreadyExists(message string) *AppError {
	return New(ErrCodeUserConflict, message, http.StatusConflict)
}
