package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"todo_list_go/internal/service"
	customErrors "todo_list_go/pkg/errors"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("sign-up", h.signUp)
		users.POST("sign-in", h.signIn)
		users.GET("me", h.getMe)
		users.PUT("me", h.updateMe)
	}
}

type signUpUserInput struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type signInUserInput struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func (h *Handler) signUp(c *gin.Context) {
	var inp signUpUserInput

	if err := c.BindJSON(&inp); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				field := strings.ToLower(fe.Field())
				out[field] = customErrors.ValidationErrorToText(fe)
			}
			newErrorsResponse(c, http.StatusBadRequest, out)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Users.SignUp(c, service.SignUpUserInput(inp)); err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, customErrors.ErrUserAlreadyExists.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
	var inp signInUserInput
	if err := c.BindJSON(&inp); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				field := strings.ToLower(fe.Field())
				out[field] = customErrors.ValidationErrorToText(fe)
			}
			newErrorsResponse(c, http.StatusBadRequest, out)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.services.Users.SignIn(c, service.SignInUserInput(inp))
	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{accessToken})
}

func (h *Handler) getMe(c *gin.Context) {

}

func (h *Handler) updateMe(c *gin.Context) {

}
