package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo_list_go/internal/service"
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

func (h *Handler) signUp(c *gin.Context) {
	var inp signUpUserInput

	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Users.SignUp(c, service.SignUpUserInput(inp)); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {

}

func (h *Handler) getMe(c *gin.Context) {

}

func (h *Handler) updateMe(c *gin.Context) {

}
