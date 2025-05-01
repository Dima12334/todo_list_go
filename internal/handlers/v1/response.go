package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo_list_go/pkg/logger"
)

type errorBodyResponse struct {
	Type    string `json:"type" binding:"oneof=string dict"`
	Details any    `json:"details"`
}

type errorResponse struct {
	Error errorBodyResponse `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, err any) {
	if statusCode == http.StatusInternalServerError {
		logger.Errorf("server error: %v", err)
		c.AbortWithStatusJSON(
			statusCode, errorResponse{
				Error: errorBodyResponse{
					Type:    "string",
					Details: "internal server error",
				},
			},
		)
		return
	}

	switch err.(type) {
	case string:
		c.AbortWithStatusJSON(statusCode, errorResponse{
			Error: errorBodyResponse{
				Type:    "string",
				Details: err,
			},
		})
	case map[string]string:
		c.AbortWithStatusJSON(statusCode, errorResponse{
			Error: errorBodyResponse{
				Type:    "dict",
				Details: err,
			},
		})
	default:
		c.AbortWithStatusJSON(statusCode, errorResponse{
			Error: errorBodyResponse{
				Type:    "string",
				Details: fmt.Sprintf("%v", err),
			},
		})
	}
}
