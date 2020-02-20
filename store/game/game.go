package trivia

import (
	"errors"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer int) error
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

func (g *impl) CreateQuiz(context ctx.CTX, user *models.User, content string, options []string, answer int) error {
	// Check input
	if len(options) < 2 {
		return errors.New("number of options should be greater than 1")
	}
	if answer < 0 || answer >= len(options) {
		return errors.New("invalid answer")
	}

	// Write db
	if _, err := g.mysql.Exec("INSERT INTO quizzes (content, option1, option2, option3, option4, answer, creator) VALUES(?, ?, ?, ?, ?, ?, ?)", content, options[0], options[1], options[2], options[3], answer, user.ID); err != nil {
		context.WithField("err", err).Error("CreateQuiz failed at mysql.Exec")
		return err
	}
	return nil
}
