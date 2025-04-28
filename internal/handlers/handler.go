package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	api_v1 "todo_list_go/internal/handlers/v1"
	"todo_list_go/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
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
	handlerAPIV1 := api_v1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerAPIV1.Init(api)
	}
}
