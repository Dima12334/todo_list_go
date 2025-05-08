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
		tasks.Use(h.UserIdentityMiddleware)
		tasks.GET("", h.GetAllTasks)
		tasks.POST("", h.CreateTask)
		tasks.GET("/:id", h.GetTaskById)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
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

// GetAllTasks @Summary Get Tasks
// @Security ApiKeyAuth
// @Tags tasks
// @Description get tasks
// @ModuleID getTasks
// @Accept  json
// @Produce  json
// @Param page query int false "page number" default(1)
// @Param limit query int false "items per page" default(20)
// @Param completed query bool false "completed (true/false)"
// @Param createdAtDateFrom query string false "format: yyyy-mm-dd"
// @Param createdAtDateTo query string false "format: yyyy-mm-dd"
// @Param categoryIds query string false "Comma-separated list of category IDs (e.g. uuid1,uuid2)"
// @Success 200 {array} taskResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /tasks [get]
func (h *Handler) GetAllTasks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var query domain.GetTasksQuery
	if err := c.BindQuery(&query); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, query)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
			return
		}
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	query.NormalizePagination()
	query.NormalizeFilters()

	res, err := h.services.Tasks.GetList(c, userID, query)
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
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: res.TotalPages,
		Total:      res.TotalItems,
		Items:      tasksList,
	})
}

// CreateTask @Summary Create Task
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
func (h *Handler) CreateTask(c *gin.Context) {
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

// GetTaskById @Summary Get Task
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
func (h *Handler) GetTaskById(c *gin.Context) {
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

// UpdateTask @Summary Update Task
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
func (h *Handler) UpdateTask(c *gin.Context) {
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

// DeleteTask @Summary Delete Task
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
func (h *Handler) DeleteTask(c *gin.Context) {
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
