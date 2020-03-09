package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/asymptoter/practice-backend/apis/auth"
	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/db"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	sql *sqlx.DB
}

func (s *AuthTestSuite) SetupTest() {
	config.Init(". ")
	s.sql = db.MustNew("postgres", false)
	_, err := s.sql.Exec("TRUNCATE users")
	s.NoError(err)
}

func (s *AuthTestSuite) TearDownTest() {
	_, err := s.sql.Exec("TRUNCATE users")
	s.NoError(err)
	s.sql.Close()
}

func (s *AuthTestSuite) TestSignup() {
	cfg := config.Value.Server.Testing
	body, _ := json.Marshal(auth.SignupRequest{
		Email:    cfg.Email.Account,
		Password: cfg.Email.Password,
	})

	resp, err := http.Post("http://127.0.0.1:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
