package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo_list_go/internal/domain"
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
	Description string `json:"description" binding:"min=0,max=255"`
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

func toCategoryResponse(category domain.Category) categoryResponse {
	return categoryResponse{
		ID:          category.ID,
		CreatedAt:   category.CreatedAt,
		Title:       category.Title,
		Description: category.Description,
		Color:       category.Color,
	}
}

// @Summary Get Tasks
// @Security ApiKeyAuth
// @Tags tasks
// @Description get tasks
// @ModuleID getTasks
// @Accept  json
// @Produce  json
// @Param page query int false "page number" default(1)
// @Param limit query int false "items per page" default(20)
// @Success 200 {array} taskResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks [get]
func (h *Handler) getAllTasks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var pagination domain.PaginationQuery
	if err := c.BindQuery(&pagination); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, pagination)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pagination.NormalizePagination()

	res, err := h.services.Tasks.GetList(c, userID, pagination)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	tasksList := make([]taskResponse, len(res.Items))
	for i, task := range res.Items {
		tasksList[i] = taskResponse{
			ID:          task.ID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Category:    toCategoryResponse(task.Category),
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		}
	}

	c.JSON(http.StatusOK, paginatedResponse[taskResponse]{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: res.TotalPages,
		Total:      res.TotalItems,
		Items:      tasksList,
	})
}

// @Summary Create Task
// @Security ApiKeyAuth
// @Tags tasks
// @Description create task
// @ModuleID createTask
// @Accept  json
// @Produce  json
// @Param input body createTaskInput true "task info"
// @Success 201 {object} taskResponse
// @Failure 400,401,409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks [post]
func (h *Handler) createTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var inp createTaskInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
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

	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrCategoryNotFound):
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, customErrors.ErrTaskAlreadyExists):
			newErrorResponse(c, http.StatusConflict, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(
		http.StatusCreated,
		taskResponse{
			ID:          task.ID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Category:    toCategoryResponse(task.Category),
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)
}

// @Summary Get Task
// @Security ApiKeyAuth
// @Tags tasks
// @Description get task
// @ModuleID getTask
// @Accept  json
// @Produce  json
// @Param id path string true "task id"
// @Success 200 {object} taskResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks/{id} [get]
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
			ID:          task.ID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Category:    toCategoryResponse(task.Category),
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)
}

// @Summary Update Task
// @Security ApiKeyAuth
// @Tags tasks
// @Description update task
// @ModuleID updateTask
// @Accept  json
// @Param id path string true "task id"
// @Param input body updateTaskInput true "update task info"
// @Success 200 {object} taskResponse
// @Failure 400,401,404,409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks/{id} [put]
func (h *Handler) updateTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var inp updateTaskInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
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
		case errors.Is(err, customErrors.ErrCategoryNotFound):
			newErrorResponse(c, http.StatusBadRequest, err.Error())
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
			ID:          task.ID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Category:    toCategoryResponse(task.Category),
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	)

}

// @Summary Delete Task
// @Security ApiKeyAuth
// @Tags tasks
// @Description delete task
// @ModuleID deleteTask
// @Accept  json
// @Produce  json
// @Param id path string true "task id"
// @Success 204
// @Failure 401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks/{id} [delete]
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
