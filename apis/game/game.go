package game

import (
	"net/http"

	"github.com/asymptoter/geochallenge-backend/apis/middleware"
	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/base/redis"
	"github.com/asymptoter/geochallenge-backend/models"
	"github.com/asymptoter/geochallenge-backend/store/game"
	"github.com/asymptoter/geochallenge-backend/store/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type handler struct {
	mysql *sqlx.DB
	redis redis.Service
	game  game.Store
}

func SetHttpHandler(r *gin.RouterGroup, db *sqlx.DB, redisService redis.Service, gs game.Store, us user.Store) {
	h := &handler{
		mysql: db,
		redis: redisService,
		game:  gs,
	}

	r.Use(middleware.GetUser(us))

	r.Handle("POST", "/quiz", h.createQuiz)
	r.Handle("GET", "/quiz", h.getQuiz)
	r.Handle("GET", "/quizzes", h.listQuizzes)
	r.Handle("DELETE", "/quiz", h.deleteQuiz)
	r.Handle("POST", "/game", h.createGame)
	r.Handle("GET", "/game", h.getGame)
	r.Handle("GET", "/games", h.listGames)
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

	if err := h.game.CreateQuiz(context, user, req.Content, req.Options, req.Answer, req.CountDown); err != nil {
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
