package trivia

import (
	"testing"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/docker"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type triviaSuite struct {
	suite.Suite
	db     *sqlx.DB
	redis  redis.Service
	trivia Store
}

func TestTriviaSuite(t *testing.T) {
	suite.Run(t, new(triviaSuite))
}

func (s *triviaSuite) SetupSuite() {
	s.db = docker.GetPostgreSQL()
	s.redis = docker.GetRedis()
	s.trivia = New(s.db, s.redis)
	s.initDB()
}

func (s *triviaSuite) initDB() {
	_, err := s.db.Exec("CREATE SEQUENCE IF NOT EXISTS quizzes_id_seq;")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS quizzes (id INT NOT NULL DEFAULT nextval('quizzes_id_seq'), content VARCHAR(512), image_url VARCHAR(100), options VARCHAR(64) ARRAY, answer VARCHAR(64), creator UUID, category VARCHAR(64), PRIMARY KEY (id));")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE INDEX IF NOT EXISTS creator_category_idx ON quizzes (creator, category);")
	s.Require().NoError(err)
	_, err = s.db.Exec("ALTER SEQUENCE quizzes_id_seq OWNED BY quizzes.id;")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS games (id UUID, name VARCHAR(32), quiz_ids INT ARRAY, mode SMALLINT, count_down SMALLINT, creator UUID, PRIMARY KEY (id));")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE INDEX IF NOT EXISTS creator_name_idx ON games (creator, name);")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS game_results (user_id UUID, game_id UUID, play_date BIGINT, correct_count INT, time_spent BIGINT, PRIMARY KEY (user_id, game_id, play_date));")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE INDEX IF NOT EXISTS game_id_idx ON game_results (game_id);")
	s.Require().NoError(err)
}

func (s *triviaSuite) SetupTest() {
}

func (s *triviaSuite) TearDownTest() {
	_, err := s.db.Exec("TRUNCATE TABLE quizzes RESTART IDENTITY;")
	s.Require().NoError(err)
	_, err = s.db.Exec("TRUNCATE TABLE games;")
	s.Require().NoError(err)
	_, err = s.db.Exec("TRUNCATE TABLE game_results;")
	s.Require().NoError(err)
}

func (s *triviaSuite) TearDownSuite() {
}

func (s *triviaSuite) TestCreateQuiz() {
	context := ctx.Background()
	userID := uuid.New()
	q := &models.Quiz{
		Content:  "quiz1",
		Options:  pq.StringArray{"1", "2", "3", "4"},
		Answer:   "4",
		Creator:  userID,
		Category: "no",
	}
	s.NoError(s.trivia.CreateQuiz(context, q))
}

func (s *triviaSuite) TestCreateGame() {
	context := ctx.Background()
	userID := uuid.New()
	q := &models.Quiz{
		Content:  "quiz1",
		Options:  pq.StringArray{"1", "2", "3", "4"},
		Answer:   "4",
		Creator:  userID,
		Category: "no",
	}
	s.NoError(s.trivia.CreateQuiz(context, q))
	g := &models.Game{
		ID:        uuid.New(),
		Name:      "game1",
		QuizIDs:   pq.Int64Array{1},
		Mode:      models.TriviaModePlayAll,
		CountDown: 10,
		Creator:   userID,
	}
	s.NoError(s.trivia.CreateGame(context, g))
}

func (s *triviaSuite) TestPlayGame() {
	context := ctx.Background()
	userID := uuid.New()

	s.NoError(s.trivia.CreateQuiz(context, &models.Quiz{
		Content:  "quiz1",
		Options:  pq.StringArray{"1", "2", "3", "4"},
		Answer:   "4",
		Creator:  userID,
		Category: "no",
	}))
	q, err := s.trivia.GetQuizzes(context, userID, "", "")
	s.NoError(err)
	s.Equal(int64(1), q[0].ID)

	s.NoError(s.trivia.CreateGame(context, &models.Game{
		Name:      "game1",
		QuizIDs:   pq.Int64Array{1},
		Mode:      models.TriviaModePlayAll,
		CountDown: 10,
		Creator:   userID,
	}))

	g, err := s.trivia.GetGames(context, userID, "game1")
	s.NoError(err)
	s.Equal(1, len(g))

	_, _, err = s.trivia.StartGame(context, userID, g[0].ID)
	s.NoError(err)

	_, res, err := s.trivia.Answer(context, userID, g[0].ID, "4")
	s.NoError(err)
	s.Equal(1, res.CorrectCount)
}
