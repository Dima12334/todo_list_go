package v1

import (
	"github.com/gin-gonic/gin"
	"todo_list_go/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initCategoriesRoutes(v1)
		h.initTasksRoutes(v1)
	}
}
