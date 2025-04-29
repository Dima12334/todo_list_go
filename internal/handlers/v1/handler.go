package v1

import (
	"github.com/gin-gonic/gin"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		v1.Use(validateIDParam)
		h.initUsersRoutes(v1)
		h.initCategoriesRoutes(v1)
		h.initTasksRoutes(v1)
	}
}
