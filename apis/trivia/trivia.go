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

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, ts trivia.Store, us user.Store) {
	h := &handler{
		mysql:  db,
		redis:  redisService,
		trivia: ts,
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
	r.Handle("DELETE", "/game", h.deleteGame)
	// Answer a quiz in a game
	r.Handle("POST", "/answer", h.answer)
}

func (h *handler) createQuiz(c *gin.Context) {
	context := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)

	quiz := &models.Quiz{}
	if err := c.ShouldBind(quiz); err != nil {
		context.WithFields(logrus.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("createQuiz failed at ShouldBind ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	quiz.Creator = user.ID
	quiz.ImageURL = "" // TODO transfer image into url
	if err := h.trivia.CreateQuiz(context, quiz); err != nil {
		context.WithFields(logrus.Fields{
			"quiz":  quiz,
			"error": err,
		}).Error("createQuiz failed at trivia.CreateQuiz")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

type GetQuizzesRequest struct {
	Content  string `json:"content"`
	Category string `json:"category"`
}

func (h *handler) getQuizzes(c *gin.Context) {
	context := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)
	context = ctx.WithValue(context, "userID", user.ID)
	var req GetQuizzesRequest
	if err := c.ShouldBind(&req); err != nil {
		context.WithField("error", err).Error("getQuizzes failed at ShouldBind ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.trivia.GetQuizzes(context, user.ID, req.Content, req.Category)
	if err != nil {
		context.WithFields(logrus.Fields{
			"params": req,
			"error":  err,
		}).Error("createQuiz failed at trivia.CreateQuiz")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *handler) deleteQuiz(c *gin.Context) {
}

func (h *handler) createGame(c *gin.Context) {
	context := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)

	game := &models.Game{}
	if err := c.ShouldBind(game); err != nil {
		context.WithFields(logrus.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("createGame failed at ShouldBind ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	game.Creator = user.ID
	if err := h.trivia.CreateGame(context, game); err != nil {
		context.WithFields(logrus.Fields{
			"game":  game,
			"error": err,
		}).Error("createGame failed at trivia.CreateGame")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func (h *handler) getGame(c *gin.Context) {
}
func (h *handler) getGames(c *gin.Context) {
}
func (h *handler) answer(c *gin.Context) {
}
func (h *handler) deleteGame(c *gin.Context) {
}
