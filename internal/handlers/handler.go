package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "todo_list_go/internal/handlers/v1"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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
	handlerV1 := v1.NewHandler()
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
