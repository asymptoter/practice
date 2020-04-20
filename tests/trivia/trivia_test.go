package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/asymptoter/practice-backend/apis/trivia"
	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/db"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type TriviaTestSuite struct {
	suite.Suite
	sql   *sqlx.DB
	redis redis.Service
	user  user.Store
	host  string
}

func (s *TriviaTestSuite) SetupTest() {
	config.Init(". ")
	s.host = "http://" + config.Value.Server.Address
	s.sql = db.MustNew("postgres", false)
	_, err := s.sql.Exec("TRUNCATE users")
	s.NoError(err)
	_, err = s.sql.Exec("TRUNCATE quizzes")
	s.NoError(err)
	_, err = s.sql.Exec("ALTER SEQUENCE quizzes_id_seq RESTART WITH 1")
	s.NoError(err)
	_, err = s.sql.Exec("TRUNCATE games")
	s.NoError(err)

	s.redis = redis.NewService(config.Value.Redis.Address)
	s.user = user.NewStore(s.sql, s.redis)
}

func (s *TriviaTestSuite) TearDownTest() {
	s.Require().NoError(s.sql.Close())
}

func (s *TriviaTestSuite) TestTriviaFlow() {
	context := ctx.Background()

	// Prepare user
	u := &models.User{
		ID:    uuid.New(),
		Token: uuid.New(),
	}
	s.Require().NoError(s.user.Create(context, u))

	// Create quiz
	quiz := models.Quiz{
		Content:  "content",
		Options:  []string{"option1", "option2", "option3", "option4"},
		Answer:   "option4",
		Category: "",
	}
	bodyByte, _ := json.Marshal(quiz)
	body := bytes.NewBuffer(bodyByte)
	req, err := http.NewRequest("POST", s.host+"/api/v1/trivia/quiz", body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", u.Token.String())
	c := &http.Client{}
	res, err := c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusCreated, res.StatusCode)

	// Get quizzes
	body = bytes.NewBuffer([]byte{})
	req, err = http.NewRequest("GET", s.host+"/api/v1/trivia/quizzes", body)
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", u.Token.String())
	c = &http.Client{}
	res, err = c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusOK, res.StatusCode)
	b, err := ioutil.ReadAll(res.Body)
	defer s.Require().NoError(res.Body.Close())
	s.NoError(err)
	quizzes := []*models.Quiz{}
	s.Require().NoError(json.Unmarshal(b, &quizzes))
	s.Len(quizzes, 1)
	quiz.ID = 1
	quiz.Creator = u.ID
	s.Equal(quiz, *quizzes[0])

	// Create game
	game := models.Game{
		Name:      "game1",
		QuizIDs:   []int64{1},
		Mode:      models.TriviaModePlayAll,
		CountDown: 10,
		Creator:   u.ID,
	}
	bodyByte, _ = json.Marshal(game)
	body = bytes.NewBuffer(bodyByte)
	req, err = http.NewRequest("POST", s.host+"/api/v1/trivia/game", body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", u.Token.String())
	c = &http.Client{}
	res, err = c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusCreated, res.StatusCode)

	// Get game
	body = bytes.NewBuffer([]byte{})
	req, err = http.NewRequest("GET", s.host+"/api/v1/trivia/games", body)
	req.Header.Add("token", u.Token.String())
	c = &http.Client{}
	res, err = c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusOK, res.StatusCode)
	b, err = ioutil.ReadAll(res.Body)
	defer s.Require().NoError(res.Body.Close())
	s.NoError(err)
	games := []*models.Game{}
	s.Require().NoError(json.Unmarshal(b, &games))
	s.Len(games, 1)
	s.Equal(game.Name, games[0].Name)
	s.Equal(game.QuizIDs, games[0].QuizIDs)

	fmt.Println(games[0])
	// Start game
	body = bytes.NewBuffer([]byte{})
	req, err = http.NewRequest("GET", s.host+"/api/v1/trivia/game?gameID="+games[0].ID.String(), body)
	req.Header.Add("token", u.Token.String())
	c = &http.Client{}
	res, err = c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusOK, res.StatusCode)
	b, err = ioutil.ReadAll(res.Body)
	defer s.Require().NoError(res.Body.Close())
	s.NoError(err)
	res1 := trivia.StartGameResponse{}
	s.Require().NoError(json.Unmarshal(b, &res1))
	s.Equal(res1.Game.Name, games[0].Name)
	s.Equal(res1.Game.QuizIDs, games[0].QuizIDs)
	s.Equal(res1.Quiz.ID, int64(games[0].QuizIDs[0]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TriviaTestSuite))
}
