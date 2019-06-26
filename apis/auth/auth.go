package auth

import (
	"log"
	"net/http"

	"github.com/asymptoter/practice/base/ctx"
	authStore "github.com/asymptoter/practice/store/auth"
	"github.com/gin-gonic/gin"
)

type handler struct {
	authService authStore.Service
}

type signupInfo struct {
	email    string
	password string
}

func SetHttpHandler(r *gin.Engine, authService authStore.Service) {
	h := &handler{
		authService: authService,
	}

	r.POST("/signup", h.signup)
	r.POST("/login", h.login)
}

func (h *handler) signup(c *gin.Context) {
	info := &authStore.SignupInfo{}
	context := ctx.NewContext(c)
	if err := c.ShouldBind(info); err != nil {
		log.Println("c.ShouldBind failed ", err)
		c.Status(http.StatusBadRequest)
		return

	}

	if err := h.authService.Signup(context, info); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	return
}

func (h *handler) login(c *gin.Context) {
}
