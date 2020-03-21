package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	s.redis = redis.NewService()
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
	bodyByte, _ := json.Marshal(trivia.CreateQuizRequest{
		Content:   "content",
		Options:   []string{"option1", "option2", "option3", "option4"},
		Answer:    "option4",
		CountDown: 10,
	})
	body := bytes.NewBuffer(bodyByte)
	req, err := http.NewRequest("POST", s.host+"/api/v1/trivia/quiz", body)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(u.Token.String())
	req.Header.Add("token", u.Token.String())
	c := &http.Client{}
	res, err := c.Do(req)
	s.NoError(err, err)
	s.Equal(http.StatusOK, res.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TriviaTestSuite))
}
