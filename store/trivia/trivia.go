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
	CreateGame(context ctx.CTX, game *models.Game) error
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
	res := []*models.Quiz{}
	query := "SELECT id, content, image_url, options, answer, creator, category FROM quizzes WHERE creator = $1 AND content LIKE '%' || $2 || '%' AND category LIKE '%' || $3 || '%'"
	if err := s.db.SelectContext(context, &res, query, userID, content, category); err != nil {
		context.WithField("err", err).Error("GetQuizzes failed at db.SelectContext")
		return nil, err
	}

	return res, nil
}

func (s *impl) CreateGame(context ctx.CTX, g *models.Game) error {
	// Write db
	res, err := s.db.ExecContext(context, "INSERT INTO games (name, quiz_ids, mode, count_down, creator) SELECT $1, $2::int[], $3, $4, $5 WHERE NOT EXISTS ((SELECT * FROM unnest($2::int[])) EXCEPT (SELECT id FROM quizzes))", g.Name, pq.Array(g.QuizIDs), g.Mode, g.CountDown, g.Creator)
	if err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": g.Creator,
		}).Error("CreateGame failed at db.ExecContext")
		return err
	}
	context.Info(res.RowsAffected())
	return nil
}
