package game

import (
	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/base/redis"
	"github.com/asymptoter/geochallenge-backend/models"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer, countDown int) error
}

type impl struct {
	mysql *sqlx.DB
	redis redis.Service
}

func NewStore(db *sqlx.DB, redisService redis.Service) Store {
	return &impl{
		mysql: db,
		redis: redisService,
	}
}

func (g *impl) CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer, countDown int) error {
	return nil
}
