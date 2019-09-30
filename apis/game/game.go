package game

import "github.com/gin-gonic/gin"

type handler struct {
	gameService gameStore.Service
}

func SetHttpHandler(r *gin.Engine, gameService gameStore.Service) {
	h := &handler{
		gameService: gameService,
	}

	r.GET("/quiz", h.quiz)
	r.POST("/answer", h.answer)
	r.GET("/result", h.result)
}

func (h *handler) quiz(c *gin.Context) {
}

func (h *handler) answer(c *gin.Context) {
}

func (h *handler) result(c *gin.Context) {
}
