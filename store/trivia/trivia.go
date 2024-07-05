package trivia

import (
	"errors"
	"strconv"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/external/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Store interface {
	CreateQuiz(context ctx.CTX, quiz *models.Quiz) error
	GetQuizzes(context ctx.CTX, userID uuid.UUID, content, category string) ([]*models.Quiz, error)
	CreateGame(context ctx.CTX, game *models.Game) error
	GetGames(context ctx.CTX, userID uuid.UUID, name string) ([]*models.Game, error)
	StartGame(context ctx.CTX, userID, gameID uuid.UUID) (*models.Game, *models.Quiz, error)
	Answer(context ctx.CTX, userID, gameID uuid.UUID, answer string) (*models.Quiz, *models.GameResult, error)
}

type impl struct {
	db    *sqlx.DB
	redis redis.Service
}

func New(db *sqlx.DB, redisService redis.Service) Store {
	return &impl{
		db:    db,
		redis: redisService,
	}
}

func (s *impl) CreateQuiz(ctx ctx.CTX, q *models.Quiz) error {
	// Check input
	if len(q.Options) < 2 {
		return errors.New("number of options should be greater than 1")
	} else if len(q.Options) > 4 {
		return errors.New("number of options should be less than 5")
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
	if _, err := s.db.ExecContext(ctx, "INSERT INTO quizzes (content, image_url, options, answer, creator, category) VALUES($1, $2, $3, $4, $5, $6)", q.Content, q.ImageURL, pq.Array(q.Options), q.Answer, q.Creator, q.Category); err != nil {
		ctx.With(
			"conent", q.Content,
			"userID", q.Creator,
		).Error(err)
		return err
	}
	return nil
}

func (s *impl) GetQuizzes(ctx ctx.CTX, userID uuid.UUID, content, category string) ([]*models.Quiz, error) {
	res := []*models.Quiz{}
	query := "SELECT id, content, image_url, options, answer, creator, category FROM quizzes WHERE creator = $1 AND content LIKE '%' || $2 || '%' AND category LIKE '%' || $3 || '%'"
	if err := s.db.SelectContext(ctx, &res, query, userID, content, category); err != nil {
		ctx.Error(err)
		return nil, err
	}

	return res, nil
}

func (s *impl) CreateGame(ctx ctx.CTX, g *models.Game) error {
	if len(g.QuizIDs) < 1 {
		return errors.New("the number of quiz should be positive")
	}
	if g.CountDown < 1 {
		return errors.New("count down should greater than zero")
	}

	g.ID = uuid.New()
	// Write db
	query := "INSERT INTO games (id, name, quiz_ids, mode, count_down, creator) SELECT $1, $2, $3::int[], $4, $5, $6 WHERE NOT EXISTS ((SELECT * FROM unnest($3::int[])) EXCEPT (SELECT id FROM quizzes))"
	if _, err := s.db.ExecContext(ctx, query, g.ID, g.Name, pq.Array(g.QuizIDs), g.Mode, g.CountDown, g.Creator); err != nil {
		ctx.With("user_id", g.Creator).Error(err)
		return err
	}
	return nil
}

func (s *impl) GetGames(ctx ctx.CTX, userID uuid.UUID, name string) ([]*models.Game, error) {
	res := []*models.Game{}
	query := "SELECT id, name, quiz_ids, mode, count_down, creator FROM games WHERE creator = $1 AND name LIKE '%' || $2 || '%'"
	if err := s.db.SelectContext(ctx, &res, query, userID, name); err != nil {
		ctx.Error(err)
		return nil, err
	}
	return res, nil
}

func (s *impl) StartGame(ctx ctx.CTX, userID, gameID uuid.UUID) (*models.Game, *models.Quiz, error) {
	g := &models.Game{}
	query := "SELECT id, name, quiz_ids, mode, count_down, creator FROM games WHERE id = $1"
	if err := s.db.GetContext(ctx, g, query, gameID); err != nil {
		ctx.Error(err)
		return nil, nil, err
	}

	quizzes := []*models.Quiz{}
	query = "SELECT id, content, image_url, options, answer FROM quizzes WHERE id IN (SELECT * FROM unnest($1::int[]))"
	if err := s.db.SelectContext(ctx, &quizzes, query, g.QuizIDs); err != nil {
		ctx.With("quiz_ids", g.QuizIDs).Error(err)
		return nil, nil, err
	}

	status := &models.GameStatus{
		QuizNo:    0,
		QuizIDs:   g.QuizIDs,
		Answers:   []string{},
		Mode:      g.Mode,
		CountDown: g.CountDown,
		StartTime: time.Now().Unix(),
	}
	for _, q := range quizzes {
		key := "trivia:quizID:" + strconv.FormatInt(q.ID, 10)
		if err := s.redis.Set(ctx, key, q, 10*time.Minute); err != nil {
			ctx.With(
				"key", key,
				"value", q,
			).Error(err)
			return nil, nil, err
		}
		status.CorrectAnswers = append(status.CorrectAnswers, q.Answer)
	}

	key := "trivia:userID:" + userID.String() + ":gameID:" + gameID.String()
	if err := s.redis.Set(ctx, key, status, 10*time.Minute); err != nil {
		ctx.With(
			"key", key,
			"value", status,
		).Error(err)
		return nil, nil, err
	}

	quizzes[0].Answer = ""
	return g, quizzes[0], nil
}

func (s *impl) Answer(ctx ctx.CTX, userID, gameID uuid.UUID, answer string) (*models.Quiz, *models.GameResult, error) {
	key := "trivia:userID:" + userID.String() + ":gameID:" + gameID.String()
	status := &models.GameStatus{}
	if err := s.redis.Get(ctx, key, status); err != nil {
		ctx.With("key", key).Error(err)
		return nil, nil, err
	}

	status.QuizNo++
	status.Answers = append(status.Answers, answer)
	if isEnd(status) {
		return nil, s.calculateGameResult(ctx, status, userID, gameID), nil
	}

	q := &models.Quiz{}
	key = "trivia:quizID:" + strconv.FormatInt(int64(status.QuizIDs[status.QuizNo]), 10)
	if err := s.redis.Get(ctx, key, q); err != nil {
		ctx.With("key", key).Error(err)
		return nil, nil, err
	}

	return q, nil, nil
}

func isEnd(s *models.GameStatus) bool {
	if s.QuizNo == len(s.QuizIDs) {
		return true
	}
	if s.Mode == models.TriviaModeNoWrong && s.Answers[s.QuizNo] != s.CorrectAnswers[s.QuizNo] {
		return true
	}
	return false
}

func (s *impl) calculateGameResult(ctx ctx.CTX, status *models.GameStatus, userID, gameID uuid.UUID) *models.GameResult {
	now := time.Now().Unix()
	res := &models.GameResult{
		UserID:    userID,
		GameID:    gameID,
		PlayDate:  now,
		TimeSpent: now - status.StartTime,
	}

	for i, answer := range status.Answers {
		if answer == status.CorrectAnswers[i] {
			res.CorrectCount++
		}
	}

	// Store game result
	query := "INSERT into game_results (user_id, game_id, play_date, correct_count, time_spent) VALUES ($1, $2, $3, $4, $5)"
	if _, err := s.db.ExecContext(ctx, query, res.UserID, res.GameID, res.PlayDate, res.CorrectCount, res.TimeSpent); err != nil {
		ctx.With("res", res).Error(err)
		return nil
	}
	return res
}
