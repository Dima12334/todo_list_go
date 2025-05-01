package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo_list_go/internal/service"
	customErrors "todo_list_go/pkg/errors"
)

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

type createCategoryInput struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"required,min=0,max=255"`
	Color       string `json:"color" binding:"oneof=red blue yellow purple green brown"`
}

type updateCategoryInput struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,min=0,max=255"`
	Color       *string `json:"color" binding:"omitempty,oneof=red blue yellow purple green brown"`
}

type categoryResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

// @Summary Get Categories
// @Security ApiKeyAuth
// @Tags categories
// @Description get categories
// @ModuleID getCategories
// @Accept  json
// @Produce  json
// @Success 200 {array} categoryResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /categories [get]
func (h *Handler) getAllCategories(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categories, err := h.services.Categories.GetList(c, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categoriesList := make([]categoryResponse, len(categories))
	for i, category := range categories {
		categoriesList[i] = categoryResponse{
			ID:          category.ID,
			CreatedAt:   category.CreatedAt,
			Title:       category.Title,
			Description: category.Description,
			Color:       category.Color,
		}
	}

	c.JSON(http.StatusOK, categoriesList)
}

// @Summary Create Category
// @Security ApiKeyAuth
// @Tags categories
// @Description create category
// @ModuleID createCategory
// @Accept  json
// @Produce  json
// @Param input body createCategoryInput true "category info"
// @Success 201 {object} categoryResponse
// @Failure 400,401,409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /categories [post]
func (h *Handler) createCategory(c *gin.Context) {
	var inp createCategoryInput

	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	category, err := h.services.Categories.Create(
		c,
		service.CreateCategoryInput{
			UserID:      userID,
			Title:       inp.Title,
			Description: inp.Description,
			Color:       inp.Color,
		},
	)
	if err != nil {
		if errors.Is(err, customErrors.ErrCategoryAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, customErrors.ErrCategoryAlreadyExists.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, categoryResponse{
		ID:          category.ID,
		CreatedAt:   category.CreatedAt,
		Title:       category.Title,
		Description: category.Description,
		Color:       category.Color,
	})
}

// @Summary Update Category
// @Security ApiKeyAuth
// @Tags categories
// @Description update category
// @ModuleID updateCategory
// @Accept  json
// @Param id path string true "category id"
// @Param input body updateCategoryInput true "update category info"
// @Success 200 {object} categoryResponse
// @Failure 400,401,404,409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /categories/{id} [put]
func (h *Handler) updateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var inp updateCategoryInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	category, err := h.services.Categories.Update(
		c,
		service.UpdateCategoryInput{
			ID:          categoryID,
			UserID:      userID,
			Title:       inp.Title,
			Description: inp.Description,
			Color:       inp.Color,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrCategoryNotFound):
			newErrorResponse(c, http.StatusNotFound, err.Error())
		case errors.Is(err, customErrors.ErrCategoryAlreadyExists):
			newErrorResponse(c, http.StatusConflict, err.Error())
		case errors.Is(err, customErrors.ErrNoUpdateFields):
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, categoryResponse{
		ID:          category.ID,
		CreatedAt:   category.CreatedAt,
		Title:       category.Title,
		Description: category.Description,
		Color:       category.Color,
	})
}

// @Summary Delete Category
// @Security ApiKeyAuth
// @Tags categories
// @Description delete category
// @ModuleID deleteCategory
// @Accept  json
// @Produce  json
// @Param id path string true "category id"
// @Success 204
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /categories/{id} [delete]
func (h *Handler) deleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Categories.Delete(c, categoryID, userID); err != nil {
		if errors.Is(err, customErrors.ErrCategoryNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
