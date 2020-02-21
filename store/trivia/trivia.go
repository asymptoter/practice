package trivia

import (
	"errors"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateQuiz(context ctx.CTX, userID, content string, options []string, answer int) error
	GetQuizzes(context ctx.CTX, userID, content string) ([]*models.Quiz, error)
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

func (s *impl) CreateQuiz(context ctx.CTX, userID, content string, options []string, answer int) error {
	// Check input
	if len(options) < 2 {
		return errors.New("number of options should be greater than 1")
	}
	if answer < 0 || answer >= len(options) {
		return errors.New("invalid answer")
	}

	// Write db
	if _, err := s.mysql.Exec("INSERT INTO quizzes (content, option1, option2, option3, option4, answer, creator) VALUES(?, ?, ?, ?, ?, ?, ?)", content, options[0], options[1], options[2], options[3], answer, userID); err != nil {
		context.WithField("err", err).Error("CreateQuiz failed at mysql.Exec")
		return err
	}
	return nil
}

func (s *impl) GetQuizzes(context ctx.CTX, userID, content string) ([]*models.Quiz, error) {
	return []*models.Quiz{}, nil
}
