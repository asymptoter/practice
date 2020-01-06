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

type Handler struct {
	DB *gorm.DB
}

func SetHttpHandler(r *gin.Engine, db *gorm.DB) {
	h := Handler{DB: db}
	r.Handle("POST", "/auth/signup", h.signup)
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
}

func (h *Handler) signup(c *gin.Context) {
	var signupInfo SignupRequest
	context := ctx.Background()
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

	userID, err := uuid.NewRandom()
	if err != nil {
		context.WithField("err", err).Error("uuid.NewRandom failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user := models.User{
		ID:       userID.String(),
		Email:    signupInfo.Email,
		Password: string(hashedPassword),
		Token:    token.String(),
	}

	//if err := h.DB.Table("Users").Create(user).Error; err != nil {
	if err := h.DB.Create(user).Error; err != nil {
		//if err := h.DB.Exec("Insert into users ('userID', 'email', 'password', 'token') values (?, ?, ?, ?)", user.UserID
		context.WithField("err", err).Error("Create failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, SignupResponse{
		UserID: userID.String(),
		Token:  token.String(),
	})
}
