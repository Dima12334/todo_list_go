package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("sign-up", h.signUp)
		users.POST("sign-in", h.signIn)
		users.GET("me", h.getMe)
		users.PUT("me", h.updateMe)
	}
}

func (h *Handler) signUp(c *gin.Context) {

}

func (h *Handler) signIn(c *gin.Context) {

}

func (h *Handler) getMe(c *gin.Context) {

}

func (h *Handler) updateMe(c *gin.Context) {

}
