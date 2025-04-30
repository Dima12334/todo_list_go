package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo_list_go/internal/service"
	customErrors "todo_list_go/pkg/errors"
)

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

type createTaskInput struct {
	CategoryID  string `json:"category_id" binding:"required,uuid"`
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"required,min=0,max=255"`
	Completed   bool   `json:"completed"`
}

type updateTaskInput struct {
	CategoryID  *string `json:"category_id" binding:"omitempty,uuid"`
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,min=0,max=255"`
	Completed   *bool   `json:"completed" binding:"omitempty"`
}

type taskResponse struct {
	ID          string           `json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Category    categoryResponse `json:"category"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Completed   bool             `json:"completed"`
}

func (h *Handler) getAllTasks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tasks, err := h.services.Tasks.GetList(c, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tasksList := make([]taskResponse, len(tasks))
	for i, task := range tasks {
		tasksList[i] = taskResponse{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
			Category: categoryResponse{
				ID:          task.Category.ID,
				CreatedAt:   task.Category.CreatedAt,
				Title:       task.Category.Title,
				Description: task.Category.Description,
				Color:       task.Category.Color,
			},
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		}
	}

	c.JSON(http.StatusOK, tasksList)
}

func (h *Handler) createTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var inp createTaskInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err)
		if out != nil {
			newErrorsResponse(c, http.StatusBadRequest, out)
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.services.Tasks.Create(c, service.CreateTaskInput{
		UserID:      userID,
		CategoryID:  inp.CategoryID,
		Title:       inp.Title,
		Description: inp.Description,
		Completed:   inp.Completed,
	})

	c.JSON(
		http.StatusCreated,
		taskResponse{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
			Category: categoryResponse{
				ID:          task.Category.ID,
				CreatedAt:   task.Category.CreatedAt,
				Title:       task.Category.Title,
				Description: task.Category.Description,
				Color:       task.Category.Color,
			},
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)
}

func (h *Handler) getTaskById(c *gin.Context) {
	taskId := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	task, err := h.services.Tasks.GetByID(c, taskId, userID)
	if err != nil {
		if errors.Is(err, customErrors.ErrTaskNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		taskResponse{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
			Category: categoryResponse{
				ID:          task.Category.ID,
				CreatedAt:   task.Category.CreatedAt,
				Title:       task.Category.Title,
				Description: task.Category.Description,
				Color:       task.Category.Color,
			},
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)
}

func (h *Handler) updateTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var inp updateTaskInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err)
		if out != nil {
			newErrorsResponse(c, http.StatusBadRequest, out)
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.services.Tasks.Update(
		c,
		service.UpdateTaskInput{
			ID:          taskID,
			UserID:      userID,
			CategoryID:  inp.CategoryID,
			Title:       inp.Title,
			Description: inp.Description,
			Completed:   inp.Completed,
		})
	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrTaskNotFound):
			newErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, customErrors.ErrTaskAlreadyExists):
			newErrorResponse(c, http.StatusConflict, err.Error())
		case errors.Is(err, customErrors.ErrNoUpdateFields):
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(
		http.StatusOK,
		taskResponse{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
			Category: categoryResponse{
				ID:          task.Category.ID,
				CreatedAt:   task.Category.CreatedAt,
				Title:       task.Category.Title,
				Description: task.Category.Description,
				Color:       task.Category.Color,
			},
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)

}

func (h *Handler) deleteTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Tasks.Delete(c, taskID, userID)
	if err != nil {
		if errors.Is(err, customErrors.ErrTaskNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
