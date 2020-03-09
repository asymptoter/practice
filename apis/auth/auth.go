package auth

import (
	"encoding/json"
	"errors"
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
	user  user.Store
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, us user.Store) {
	h := handler{
		sql:   db,
		redis: redisService,
		user:  us,
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
	UserID string `json:"userID"`
	Token  string `json:"token"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("bcrypt.GenerateFromPassword failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	token := uuid.New().String()
	userID := uuid.New().String()
	activeToken := uuid.New().String()

	user := models.User{
		ID:       userID,
		Email:    signupInfo.Email,
		Password: string(hashedPassword),
		Token:    token,
	}

	if _, err := h.sql.Exec("INSERT INTO users (id, email, password, register_date) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.Password, time.Now().Unix()); err != nil {
		context.WithField("err", err).Error("Create failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	activeAccountKey := "auth:active:activeToken:" + activeToken
	if err := h.redis.Set(context, activeAccountKey, token, 24*time.Hour); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": userID,
		}).Error("redis.Set failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	cfg := config.Value.Server
	query := "http://" + cfg.Address + "/api/v1/auth/activation?id" + userID + "&activeToken=" + activeToken
	activeMessage := "<p>Thank you for registering at demo site.</p><p>To activate your account, please click on this link: <a href='" + query + "'>Here</a></p><p>Regards Site Admin</p>"

	if err := email.Send(context, signupInfo.Email, activeMessage); err != nil {
		context.WithField("err", err).Error("email.Send failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SignupResponse{
		UserID: userID,
	})
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

	if _, err := h.sql.Exec("UPDATE users SET token=? WHERE id=?;", string(token), userID); err != nil {
		context.WithField("err", err).Error("sql.Update failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
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
	if err := h.sql.Get(&user, "SELECT password FROM users where email = ?;", user.Email); err != nil {
		context.WithField("err", err).Error("Where failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	context = ctx.WithValue(context, "userID", user.ID)

	if string(hashedPassword) != user.Password {
		context.Error("Invalid email or password")
		c.JSON(http.StatusUnauthorized, errors.New("Invalid email or password"))
		return
	}

	userInfoKey := "user:" + user.ID
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

	c.JSON(http.StatusOK, user)
}
