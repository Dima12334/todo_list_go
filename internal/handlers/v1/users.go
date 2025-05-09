package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo_list_go/internal/service"
	customErrors "todo_list_go/pkg/errors"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("sign-up", h.SignUp)
		users.POST("sign-in", h.SignIn)
		authenticated := users.Group("/", h.UserIdentityMiddleware)
		{
			authenticated.GET("me", h.GetMe)
		}
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

type userMeResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

// SignUp @Summary SignUp
// @Tags users
// @Description create user
// @ModuleID createUser
// @Accept  json
// @Produce  json
// @Param input body signUpUserInput true "user info"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var inp signUpUserInput

	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
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

// SignIn @Summary SignIn
// @Tags users
// @Description login user
// @ModuleID loginUser
// @Accept  json
// @Produce  json
// @Param input body signInUserInput true "user credentials"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var inp signInUserInput
	if err := c.BindJSON(&inp); err != nil {
		out := customErrors.FormatValidationErrorOutput(err, inp)
		if out != nil {
			newErrorResponse(c, http.StatusBadRequest, out)
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

// GetMe @Summary Get me
// @Security ApiKeyAuth
// @Tags users
// @Description get current user
// @ModuleID getMe
// @Accept  json
// @Produce  json
// @Success 200 {object} userMeResponse
// @Failure 401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := h.services.Users.GetByID(c, userID)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userMeResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Email:     user.Email,
	})
}
