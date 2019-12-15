package auth

import (
	"net/http"

	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func SetHttpHandler(r *gin.Engine, db *gorm.DB) {
	h := &AuthHandler{DB: db}
	r.POST("/auth/signup", h.signup)
}

type AuthHandler struct {
	DB *gorm.DB
}

type SignupRequest struct {
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
}

type SignupReply struct {
	Token string `json:"token"`
}

func (h *AuthHandler) signup(c *gin.Context) {
	var signupInfo SignupRequest
	context := ctx.NewContext(c)
	if err := c.ShouldBind(&signupInfo); err != nil {
		context.WithFields(logrus.Fields{
			"params": signupInfo,
			"error":  err,
		}).Error("c.ShouldBind failed")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("bcrypt.GenerateFromPassword failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		context.WithField("err", err).Error("uuid.NewRandom failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user := models.User{
		Email:    signupInfo.Email,
		Password: hashedPassword,
		Token:    token.String(),
	}

	h.DB.Table("users").Create(user)
}
