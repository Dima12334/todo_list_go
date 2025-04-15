package v1

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initCategoriesRoutes(v1)
		h.initTasksRoutes(v1)
	}
}
