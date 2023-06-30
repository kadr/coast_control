package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func (r *Response) Success(c *gin.Context, statusCode int, data any) *Response {
	if data == nil {
		r.Data = map[string]bool{"success": true}

	} else {
		r.Data = data
	}
	c.JSON(statusCode, r.Data)

	return &Response{Data: data}
}

func (r *Response) Error(c *gin.Context, statusCode int, message string) *Response {
	r.Message = message
	c.AbortWithStatusJSON(statusCode, r.Message)

	return &Response{Message: message}
}
