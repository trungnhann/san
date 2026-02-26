package response

import (
	"net/http"

	"san/pkg/apperr"

	"github.com/gin-gonic/gin"
)

// Envelope is the standard API response structure
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

// ErrorData holds error details
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MetaData holds pagination or other metadata
type MetaData struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Success sends a success response with data
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a success response with data and metadata
func SuccessWithMeta(c *gin.Context, status int, data interface{}, meta MetaData) {
	c.JSON(status, Envelope{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	var appErr *apperr.AppError
	var status int
	var errData ErrorData

	if e, ok := err.(*apperr.AppError); ok {
		appErr = e
		status = appErr.HTTPStatus
		errData = ErrorData{
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	} else {
		// Default to 500 Internal Server Error for unknown errors
		status = http.StatusInternalServerError
		errData = ErrorData{
			Code:    "INTERNAL_ERROR",
			Message: "Internal Server Error",
		}
	}

	c.JSON(status, Envelope{
		Success: false,
		Error:   &errData,
	})
}
