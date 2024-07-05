package trivia

import (
	"net/http"

	"github.com/asymptoter/practice-backend/apis/middleware"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/external/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/trivia"
	"github.com/asymptoter/practice-backend/store/user"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type handler struct {
	mysql  *sqlx.DB
	redis  redis.Service
	trivia trivia.Store
}

func SetHttpHandler(r *gin.RouterGroup, ts trivia.Store, us user.Store) {
	h := &handler{
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
	// List games created by creator
	r.Handle("GET", "/games", h.getGames)
	// Delete game created by creator
	r.Handle("DELETE", "/game", h.deleteGame)
	// Play a game
	r.Handle("GET", "/game", h.startGame)
	// Answer a quiz in a game
	r.Handle("POST", "/answer", h.answer)
}

func (h *handler) createQuiz(c *gin.Context) {
	user := c.MustGet("userInfo").(*models.User)
	ctx := ctx.Background().With("userID", user.ID)

	quiz := &models.Quiz{}
	if err := c.ShouldBind(quiz); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ctx = ctx.With("quiz", quiz)

	quiz.Creator = user.ID
	quiz.ImageURL = "" // TODO transfer image into url
	if err := h.trivia.CreateQuiz(ctx, quiz); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

type GetQuizzesParams struct {
	Content  string `form:"content"`
	Category string `form:"category"`
}

func (h *handler) getQuizzes(c *gin.Context) {
	ctx := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)
	ctx = ctx.With("userID", user.ID)
	var req GetQuizzesParams
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.trivia.GetQuizzes(ctx, user.ID, req.Content, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *handler) deleteQuiz(c *gin.Context) {
}

func (h *handler) createGame(c *gin.Context) {
	ctx := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)

	game := &models.Game{}
	if err := c.ShouldBind(game); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	game.Creator = user.ID
	if err := h.trivia.CreateGame(ctx, game); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

type GetGamesParams struct {
	Name string `form:"name"`
}

func (h *handler) getGames(c *gin.Context) {
	ctx := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)
	ctx = ctx.With("userID", user.ID)
	var req GetGamesParams
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.trivia.GetGames(ctx, user.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

type StartGameParams struct {
	GameID string `form:"gameID"`
}

type StartGameResponse struct {
	Game *models.Game `json:"game"`
	Quiz *models.Quiz `json:"quiz"`
}

func (h *handler) startGame(c *gin.Context) {
	ctx := ctx.Background()
	user := c.MustGet("userInfo").(*models.User)
	ctx = ctx.With("userID", user.ID)

	var req StartGameParams
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	gameID := uuid.MustParse(req.GameID)
	ctx = ctx.With("gameID", gameID)

	game, quiz, err := h.trivia.StartGame(ctx, user.ID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, StartGameResponse{
		Game: game,
		Quiz: quiz,
	})
}

type AnswerRequest struct {
	GameID uuid.UUID `json:"gameID"`
	Answer string    `json:"answer"`
}

type AnswerResponse struct {
	Quiz   *models.Quiz       `json:"quiz"`
	Result *models.GameResult `json:"result"`
}

func (h *handler) answer(c *gin.Context) {
	user := c.MustGet("userInfo").(*models.User)
	ctx := ctx.Background().With("userID", user.ID)

	var req AnswerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ctx = ctx.With("gameID", req.GameID)

	quiz, res, err := h.trivia.Answer(ctx, user.ID, req.GameID, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, AnswerResponse{
		Quiz:   quiz,
		Result: res,
	})
}
func (h *handler) deleteGame(c *gin.Context) {
}
