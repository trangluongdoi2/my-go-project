package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}

func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, message)
}

// AbortWithError stops the chain and writes an error response
func AbortWithError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}
