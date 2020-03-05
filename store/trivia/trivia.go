package trivia

import (
	"errors"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	CreateQuiz(context ctx.CTX, userID, content, answer string, options []string) error
	GetQuizzes(context ctx.CTX, userID, content string) ([]*models.Quiz, error)
}

type impl struct {
	db    *sqlx.DB
	redis redis.Service
}

func NewStore(db *sqlx.DB, redisService redis.Service) Store {
	return &impl{
		db:    db,
		redis: redisService,
	}
}

func (s *impl) CreateQuiz(context ctx.CTX, userID, content, answer string, options []string) error {
	// Check input
	if len(options) < 2 {
		return errors.New("number of options should be greater than 1")
	}
	flag := true
	for _, v := range options {
		if v == answer {
			flag = false
			break
		}
	}
	if flag {
		return errors.New("answer should be included in options")
	}

	// Write db
	if _, err := s.db.ExecContext(context, "INSERT INTO quizzes (content, options, answer, creator) VALUES(?, ?, ?, ?)", content, options, answer, userID); err != nil {
		context.WithField("err", err).Error("CreateQuiz failed at db.Exec")
		return err
	}
	return nil
}

func (s *impl) GetQuizzes(context ctx.CTX, userID, content string) ([]*models.Quiz, error) {
	res := []*models.Quiz{}
	//s.db.SelectContext(context, &res, "SELECT (content, options, answer")
	return res, nil
}
