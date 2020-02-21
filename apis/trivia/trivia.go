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

	// Create a quiz
	r.Handle("POST", "/quiz", h.createQuiz)
	// List quizzes created by creator
	r.Handle("GET", "/quizzes", h.getQuizzes)
	// Delete a quiz created by creator
	r.Handle("DELETE", "/quiz", h.deleteQuiz)
	// Create a game
	r.Handle("POST", "/game", h.createGame)
	// Play a game
	r.Handle("GET", "/game", h.getGame)
	// List games created by creator
	r.Handle("GET", "/games", h.getGames)
	// Delete game created by creator
	r.Handle("DELETE", "game", h.deleteGame)
	// Answer a quiz in a game
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

	if err := h.trivia.CreateQuiz(context, user.ID, req.Content, req.Options, req.Answer); err != nil {
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

func (h *handler) getQuizzes(c *gin.Context) {
}
func (h *handler) deleteQuiz(c *gin.Context) {
}
func (h *handler) createGame(c *gin.Context) {
}
func (h *handler) getGame(c *gin.Context) {
}
func (h *handler) getGames(c *gin.Context) {
}
func (h *handler) answer(c *gin.Context) {
}
