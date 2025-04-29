package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initTasksRoutes(api *gin.RouterGroup) {
	tasks := api.Group("/tasks")
	{
		tasks.Use(h.userIdentity)
		tasks.GET("", h.getAllTasks)
		tasks.POST("", h.createTask)
		tasks.GET("/:id", h.getTaskById)
		tasks.PUT("/:id", h.updateTask)
		tasks.DELETE("/:id", h.deleteTask)
	}
}

func (h *Handler) getAllTasks(c *gin.Context) {
}

func (h *Handler) createTask(c *gin.Context) {
}

func (h *Handler) getTaskById(c *gin.Context) {
}

func (h *Handler) updateTask(c *gin.Context) {
}

func (h *Handler) deleteTask(c *gin.Context) {
}
