package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	idParamCtx          = "id"
)

func getUserID(c *gin.Context) (string, error) {
	userID, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	id, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user id type")
	}

	return id, nil
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(authorizationHeader)
	if authHeader == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	return h.tokenManager.ParseJWT(headerParts[1])
}

func (h *Handler) UserIdentityMiddleware(c *gin.Context) {
	userID, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userID)
}

func ValidateIDParamMiddleware(c *gin.Context) {
	for _, param := range c.Params {
		if strings.Contains(param.Key, idParamCtx) {
			if _, err := uuid.Parse(param.Value); err != nil {
				newErrorResponse(
					c,
					http.StatusBadRequest,
					fmt.Sprintf("invalid UUID format for url parameter %s", param.Key),
				)
			}
		}
		c.Next()
	}
}
