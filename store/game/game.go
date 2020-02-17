package game

import (
	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/models"

	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer, countDown int) error
}

type impl struct {
	mysql *sqlx.DB
	redis *redis.Client
}

func NewStore(db *sqlx.DB, redis *redis.Client) Store {
	return &impl{
		mysql: db,
		redis: redis,
	}
}

func (g *impl) CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer, countDown int) error {
	return nil
}
