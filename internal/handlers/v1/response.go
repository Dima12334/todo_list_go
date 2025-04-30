package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo_list_go/pkg/logger"
)

type errorResponse struct {
	Error string `json:"error"`
}

type errorsResponse struct {
	Errors any `json:"errors"`
}

func newErrorResponse(c *gin.Context, statusCode int, error string) {
	if statusCode == http.StatusInternalServerError {
		logger.Error(error)
		c.AbortWithStatusJSON(statusCode, errorResponse{"internal server error"})
		return
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{error})
}

func newErrorsResponse(c *gin.Context, statusCode int, errors any) {
	c.AbortWithStatusJSON(statusCode, errorsResponse{errors})
}
