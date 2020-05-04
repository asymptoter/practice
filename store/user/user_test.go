package user

import (
	"fmt"
	"testing"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/docker"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type userSuite struct {
	suite.Suite
	db        *sqlx.DB
	redis     redis.Service
	userStore Store
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (s *userSuite) SetupSuite() {
	s.db = docker.GetPostgreSQL()
	s.redis = docker.GetRedis()
	s.userStore = New(s.db, s.redis)
	s.initDB()
}

func (s *userSuite) initDB() {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS users (id UUID, token UUID, name VARCHAR(320), email VARCHAR(320) UNIQUE NOT NULL, password CHAR(60) NOT NULL, register_date BIGINT, PRIMARY KEY (id));")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE UNIQUE INDEX ON users (email);")
	s.Require().NoError(err)
	_, err = s.db.Exec("CREATE UNIQUE INDEX ON users (token);")
	s.Require().NoError(err)
}

func (s *userSuite) SetupTest() {
}

func (s *userSuite) TearDownTest() {
	_, err := s.db.Exec("TRUNCATE TABLE users")
	s.Require().NoError(err)
}

func (s *userSuite) TearDownSuite() {
}

func (s *userSuite) TestCreate() {
	context := ctx.Background()
	u := &models.User{
		Email: "a@b",
	}
	_, err := s.userStore.Create(context, u)
	s.NoError(err)
}

func (s *userSuite) TestCreateDuplicatedEmail() {
	context := ctx.Background()
	u := &models.User{
		Email: "a@b",
	}
	_, err := s.userStore.Create(context, u)
	s.NoError(err)
	_, err = s.userStore.Create(context, u)
	s.Error(err)
	fmt.Println(err)
}
