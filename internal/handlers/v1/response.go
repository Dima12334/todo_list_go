package v1

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

type errorsResponse struct {
	Errors any `json:"errors"`
}

func newErrorResponse(c *gin.Context, statusCode int, error string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{error})
}

func newErrorsResponse(c *gin.Context, statusCode int, errors any) {
	c.AbortWithStatusJSON(statusCode, errorsResponse{errors})
}
