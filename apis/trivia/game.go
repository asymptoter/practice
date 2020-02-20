package trivia

import (
	"net/http"

	"github.com/asymptoter/practice-backend/apis/middleware"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/trivia"
	"github.com/asymptoter/practice-backend/store/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type handler struct {
	mysql  *sqlx.DB
	redis  redis.Service
	trivia trivia.Store
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, gs trivia.Store, us user.Store) {
	h := &handler{
		mysql:  db,
		redis:  redisService,
		trivia: gs,
	}

	r.Use(middleware.GetUser(us))

	r.Handle("POST", "/quiz", h.createQuiz)
	r.Handle("GET", "/quiz", h.getQuiz)
	r.Handle("GET", "/quizzes", h.listQuizzes)
	r.Handle("DELETE", "/quiz", h.deleteQuiz)
	r.Handle("POST", "/trivia", h.createGame)
	r.Handle("GET", "/trivia", h.getGame)
	r.Handle("GET", "/trivias", h.listGames)
	r.Handle("POST", "/answer", h.answer)
}

type CreateQuizRequest struct {
	Content   string
	Options   []string
	Answer    int
	CountDown int
}

func (h *handler) createQuiz(c *gin.Context) {
	user := c.MustGet("userInfo").(*models.User)

	var req CreateQuizRequest
	context := ctx.Background()
	if err := c.ShouldBind(&req); err != nil {
		context.WithFields(logrus.Fields{
			"params": req,
			"userID": user.ID,
			"error":  err,
		}).Error("createQuiz failed at ShouldBind ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.trivia.CreateQuiz(context, user, req.Content, req.Options, req.Answer); err != nil {
		context.WithFields(logrus.Fields{
			"params": req,
			"userID": user.ID,
			"error":  err,
		}).Error("createQuiz failed at CreateQuiz")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *handler) getQuiz(c *gin.Context) {
}

func (h *handler) listQuizzes(c *gin.Context) {
}
func (h *handler) deleteQuiz(c *gin.Context) {
}
func (h *handler) createGame(c *gin.Context) {
}
func (h *handler) getGame(c *gin.Context) {
}
func (h *handler) listGames(c *gin.Context) {
}
func (h *handler) answer(c *gin.Context) {
}
