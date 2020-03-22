package trivia

import (
	"errors"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Store interface {
	CreateQuiz(context ctx.CTX, quiz *models.Quiz) error
	GetQuizzes(context ctx.CTX, userID uuid.UUID, content, category string) ([]*models.Quiz, error)
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

func (s *impl) CreateQuiz(context ctx.CTX, q *models.Quiz) error {
	// Check input
	if len(q.Options) < 2 {
		return errors.New("number of options should be greater than 1")
	}
	flag := true
	for _, v := range q.Options {
		if v == q.Answer {
			flag = false
			break
		}
	}
	if flag {
		return errors.New("answer should be included in options")
	}

	// Write db
	if _, err := s.db.ExecContext(context, "INSERT INTO quizzes (content, image_url, options, answer, creator, category) VALUES($1, $2, $3, $4, $5, $6)", q.Content, q.ImageURL, pq.Array(q.Options), q.Answer, q.Creator, q.Category); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"conent": q.Content,
			"userID": q.Creator,
		}).Error("CreateQuiz failed at db.ExecContext")
		return err
	}
	return nil
}

func (s *impl) GetQuizzes(context ctx.CTX, userID uuid.UUID, content, category string) ([]*models.Quiz, error) {
	rows := &sqlx.Rows{}
	var err error

	query := "SELECT id, content, image_url, options, answer, creator, category FROM quizzes WHERE creator = $1 AND content LIKE '%' || $2 || '%'"
	if len(category) > 0 {
		query = query + " AND category = $3"
		rows, err = s.db.QueryxContext(context, query, userID, content, category)
	} else {
		rows, err = s.db.QueryxContext(context, query, userID, content)
	}
	if err != nil {
		context.WithField("err", err).Error("GetQuizzes failed at db.QueryxContext")
		return nil, err
	}

	res := []*models.Quiz{}
	for rows.Next() {
		t := &models.Quiz{}
		if err := rows.Scan(&t.ID, &t.Content, &t.ImageURL, (*pq.StringArray)(&t.Options), &t.Answer, &t.Creator, &t.Category); err != nil {
			context.WithField("err", err).Error("GetQuizzes failed at rows.Scan")
		}
		res = append(res, t)
	}

	return res, nil
}
