package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initCategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		categories.Use(h.userIdentity)
		categories.GET("", h.getAllCategories)
		categories.POST("", h.createCategory)
		categories.PUT("/:id", h.updateCategory)
		categories.DELETE("/:id", h.deleteCategory)
	}
}

func (h *Handler) getAllCategories(c *gin.Context) {
}

func (h *Handler) createCategory(c *gin.Context) {
}

func (h *Handler) updateCategory(c *gin.Context) {
}

func (h *Handler) deleteCategory(c *gin.Context) {
}
