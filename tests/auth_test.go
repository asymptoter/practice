package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/asymptoter/geochallenge-backend/apis/auth"
	"github.com/google/uuid"

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
	body, _ := json.Marshal(auth.SignupRequest{
		Email:    "test" + uuid.New().String(),
		Password: "test",
	})

	resp, err := http.Post("http://127.0.0.1:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
