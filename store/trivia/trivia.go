package trivia

import (
	"errors"
	"strconv"
	"time"

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
	GetGames(context ctx.CTX, userID uuid.UUID, name string) ([]*models.Game, error)
	StartGame(context ctx.CTX, userID, gameID uuid.UUID) (*models.Game, *models.Quiz, error)
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
	if len(g.QuizIDs) < 1 {
		return errors.New("the number of quiz should be positive")
	}
	if g.CountDown < 1 {
		return errors.New("count down should greater than zero")
	}

	g.ID = uuid.New()
	// Write db
	query := "INSERT INTO games (id, name, quiz_ids, mode, count_down, creator) SELECT $1, $2, $3::int[], $4, $5, $6 WHERE NOT EXISTS ((SELECT * FROM unnest($3::int[])) EXCEPT (SELECT id FROM quizzes))"
	if _, err := s.db.ExecContext(context, query, g.ID, g.Name, pq.Array(g.QuizIDs), g.Mode, g.CountDown, g.Creator); err != nil {
		context.WithFields(logrus.Fields{
			"err":    err,
			"userID": g.Creator,
		}).Error("CreateGame failed at db.ExecContext")
		return err
	}
	return nil
}

func (s *impl) GetGames(context ctx.CTX, userID uuid.UUID, name string) ([]*models.Game, error) {
	res := []*models.Game{}
	query := "SELECT id, name, quiz_ids, mode, count_down, creator FROM games WHERE creator = $1 AND name LIKE '%' || $2 || '%'"
	if err := s.db.SelectContext(context, &res, query, userID, name); err != nil {
		context.WithField("err", err).Error("GetGames failed at db.SelectContext")
		return nil, err
	}
	return res, nil
}

func (s *impl) StartGame(context ctx.CTX, userID, gameID uuid.UUID) (*models.Game, *models.Quiz, error) {
	g := &models.Game{}
	query := "SELECT id, name, quiz_ids, mode, count_down, creator FROM games WHERE id = $1"
	if err := s.db.GetContext(context, g, query, gameID); err != nil {
		context.WithField("err", err).Error("StartGame failed at db.GetContext")
		return nil, nil, err
	}

	quizzes := []*models.Quiz{}
	query = "SELECT id, content, image_url, options, answer FROM quizzes WHERE id IN (SELECT * FROM unnest($1::int[]))"
	if err := s.db.SelectContext(context, &quizzes, query, g.QuizIDs); err != nil {
		context.WithFields(logrus.Fields{
			"err":     err,
			"quizIDs": g.QuizIDs,
		}).Error("StartGame failed at db.SelectContext")
		return nil, nil, err
	}

	for _, q := range quizzes {
		key := "trivia:quizID:" + strconv.FormatInt(q.ID, 10)
		if err := s.redis.Set(context, key, q, 10*time.Minute); err != nil {
			context.WithFields(logrus.Fields{
				"err":   err,
				"key":   key,
				"value": q,
			}).Error("StartGame failed at redis.Set")
			return nil, nil, err
		}
	}

	status := &models.GameStatus{
		Name:      g.Name,
		QuizNo:    0,
		QuizIDs:   g.QuizIDs,
		Answers:   []string{},
		Mode:      g.Mode,
		CountDown: g.CountDown,
	}
	key := "trivia:userID:" + userID.String() + ":gameID:" + gameID.String()
	if err := s.redis.Set(context, key, status, 10*time.Minute); err != nil {
		context.WithFields(logrus.Fields{
			"err":   err,
			"key":   key,
			"value": status,
		}).Error("StartGame failed at redis.Set")
		return nil, nil, err
	}

	quizzes[0].Answer = ""
	return g, quizzes[0], nil
}
