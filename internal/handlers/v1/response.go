package v1

import (
	"github.com/gin-gonic/gin"
	"todo_list_go/pkg/logger"
)

type response struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
