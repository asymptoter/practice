package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/email"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/user"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	sql   *sqlx.DB
	redis redis.Service
	us    user.Store
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, us user.Store) {
	h := handler{
		sql:   db,
		redis: redisService,
		us:    us,
	}

	r.Handle("POST", "/signup", h.signup)
	r.Handle("GET", "/activation", h.activation)
	r.Handle("POST", "/login", h.login)
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	UserID uuid.UUID `json:"userID"`
}

func (h *handler) signup(c *gin.Context) {
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

	user := &models.User{
		Email:    signupInfo.Email,
		Password: signupInfo.Password,
	}

	if err := h.us.Create(context, user); err != nil {
		context.WithFields(logrus.Fields{
			"err":   err,
			"email": user.Email,
		}).Error("signup failed at us.Create")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	activeToken := uuid.New().String()
	activeAccountKey := "auth:active:activeToken:" + activeToken
	if err := h.redis.Set(context, activeAccountKey, user.Token, 24*time.Hour); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": user.ID,
		}).Error("redis.Set failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	cfg := config.Value.Server
	link := "http://" + cfg.Address + "/api/v1/auth/activation?id=" + user.ID.String() + "&activeToken=" + activeToken
	activeMessage := fmt.Sprintf(cfg.Email.ActivationMessage, link)

	if err := email.Send(context, signupInfo.Email, activeMessage); err != nil {
		context.WithField("err", err).Error("email.Send failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		UserID: user.ID,
	})
}

type ActivationResponse struct {
	Token string `json:"token"`
}

func (h *handler) activation(c *gin.Context) {
	userID := c.Query("id")
	activeToken := c.Query("activeToken")
	context := ctx.Background()

	activeAccountKey := "auth:active:activeToken:" + activeToken
	token, err := h.redis.Get(context, activeAccountKey)
	if err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": userID,
		}).Error("redis.Get failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := h.sql.Exec("UPDATE users SET token=$1 WHERE id=$2;", string(token), userID); err != nil {
		context.WithField("err", err).Error("sql.Update failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	context.Info("sql.Exec done.")

	c.JSON(http.StatusOK, ActivationResponse{
		Token: string(token),
	})
}

type loginRequest struct {
	Email    string
	Password string
}

func (h *handler) login(c *gin.Context) {
	var loginInfo loginRequest
	context := ctx.Background()
	if err := c.ShouldBind(&loginInfo); err != nil {
		context.WithFields(logrus.Fields{
			"params": loginInfo,
			"error":  err,
		}).Error("c.ShouldBind failed")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	context = ctx.WithValue(context, "email", loginInfo.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("bcrypt.GenerateFromPassword failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user := models.User{}
	if err := h.sql.Get(&user, "SELECT ID, email, password, token, register_date FROM users WHERE email = $1 LIMIT 1;", loginInfo.Email); err != nil {
		context.WithField("err", err).Error("login failed at sql.Get")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	context = ctx.WithValue(context, "userID", user.ID)

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginInfo.Password)); err != nil {
		context.WithField("err", err).Error("Invalid email or password")
		c.JSON(http.StatusUnauthorized, errors.New("Invalid email or password"))
		return
	}

	if len(user.Token) == 0 {
		context.Error("Account is not activated")
		c.JSON(http.StatusBadRequest, errors.New("Account is not activated"))
		return
	}

	userInfoKey := "user:" + user.ID.String()
	b, err := json.Marshal(user)
	if err != nil {
		context.WithField("err", err).Error("json.Marshal failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := h.redis.Set(context, userInfoKey, string(b), 24*time.Hour); err != nil {
		context.WithField("err", err).Error("redis.Set failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
