package user

import (
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
}

func (s *userSuite) initDB() {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS users (id UUID, token UUID UNIQUE, name VARCHAR(320), email VARCHAR(320) UNIQUE NOT NULL, password CHAR(60) NOT NULL, register_date BIGINT NOT NULL, PRIMARY KEY (id));")
	s.Require().NoError(err)
}

func (s *userSuite) SetupTest() {
	s.initDB()
}

func (s *userSuite) TearDownTest() {
	_, err := s.db.Exec("DROP TABLE IF EXISTS users CASCADE")
	s.Require().NoError(err)
}

func (s *userSuite) TearDownSuite() {
	/*
		docker.Purge("postgres")
		docker.Purge("redis")
	*/
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
		Email: "c@b",
	}
	_, err := s.userStore.Create(context, u)
	s.NoError(err)
	_, err = s.userStore.Create(context, u)
	s.Error(err)
}
