package auth

import (
	"log"
	"net/http"

	authStore "github.com/asymptoter/practice/app/server/store/auth"
	"github.com/asymptoter/practice/base/ctx"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authService *authStore.Service
}

type signupInfo struct {
	email    string
	password string
}

func SetHttpHandler(r *gin.Engine, authService *authStore.Service) {
	h := &authHandler{
		auth: authService,
	}

	r.POST("/signup", h.signup)
	r.POST("/login", h.login)
}

func (h *handler) signup(c *gin.Context) {
	info := &signupInfo{}
	context := ctx.NewContext(c)
	if err := c.ShouldBind(info); err != nil {
		log.Println("c.ShouldBind failed ", err)
		c.JSON(http.StatusBadRequest)
		return

	}

	if err := h.authService.Signup(context, info); err != nil {
		c.JSON()
		return
	}

	c.JSON(http.StatusOK)
	return
}
