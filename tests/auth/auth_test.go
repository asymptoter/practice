package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/asymptoter/practice-backend/apis/auth"
	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/db"
	"github.com/asymptoter/practice-backend/base/email"

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
}

func (s *AuthTestSuite) TestAuthFlow() {
	context := ctx.Background()
	cfg := config.Value.Server.Testing.Email
	body, _ := json.Marshal(auth.SignupRequest{
		Email:    cfg.Account,
		Password: cfg.Password,
	})

	// Signup
	resp, err := http.Post("http://127.0.0.1/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode, "Signup")

	msg, err := email.Receive(context, cfg.Account, cfg.Password)
	s.NoError(err)

	msg = strings.Replace(msg, "=3D", "=", -1)
	msg = strings.Replace(msg, "=\r\n", "", -1)
	context.Info(msg)

	r, _ := regexp.Compile("http.*(H){1}")
	m := string(r.Find([]byte(msg)))
	url := m[:len(m)-2]
	context.Info("Active url:", url)
	// Active account
	resp, err = http.Get(url)
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode, "Active")
	context.Info("Active OK")

	// Login
	resp, err = http.Post("http://127.0.0.1/api/v1/auth/login", "application/json", bytes.NewBuffer(body))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode, "Login")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
