package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSONResponse(c *gin.Context, status int, success bool, message string, data any) {
	c.JSON(status, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func JSONErrorResponse(c *gin.Context, statusCode int, err error) {
	if err == nil {
		err = gin.Error{Err: errors.New("Unknown error")}
	}

	c.JSON(statusCode, gin.H{"error": err.Error()})
}
