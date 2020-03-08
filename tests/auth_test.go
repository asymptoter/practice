package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/asymptoter/practice-backend/apis/auth"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func (s *AuthTestSuite) SetupTest() {
}

func (s *AuthTestSuite) TearDownTest() {
}

func (s *AuthTestSuite) TestSignup() {
	cfg := config.Server.Testing
	body, _ := json.Marshal(auth.SignupRequest{
		Email:    cfg.Email.Account,
		Password: cfg.Email.Password,
	})

	resp, err := http.Post("http://127.0.0.1:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
