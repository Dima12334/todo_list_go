package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	apiV1 "todo_list_go/internal/handlers/v1"
	"todo_list_go/internal/service"
	"todo_list_go/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{services: services, tokenManager: tokenManager}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware)

	// Init router
	router.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initApi(router)

	return router
}

func (h *Handler) initApi(router *gin.Engine) {
	handlerAPIV1 := apiV1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerAPIV1.Init(api)
	}
}
