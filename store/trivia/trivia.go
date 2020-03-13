package trivia

import (
	"errors"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/sirupsen/logrus"

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
	if _, err := s.db.ExecContext(context, "INSERT INTO quizzes (content, options, answer, creator) VALUES($1, $2, $3, $4)", content, options, answer, userID); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"conent": content,
			"userID": userID,
		}).Error("CreateQuiz failed at db.ExecContext")
		return err
	}
	return nil
}

func (s *impl) GetQuizzes(context ctx.CTX, userID, content string) ([]*models.Quiz, error) {
	res := []*models.Quiz{}
	if err := s.db.SelectContext(context, &res, "SELECT id, content, image_url, options, answer FROM quizzes WHERE creator = $1 AND content LIKE '%$2%';", userID, content); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"conent": content,
			"userID": userID,
		}).Error("GetQuizzes failed at db.SelectContext")
		return nil, err
	}
	return res, nil
}
