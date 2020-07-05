package mongo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoSuite struct {
	suite.Suite
	mongo *mgo.Session
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(mongoSuite))
}

func (s *mongoSuite) SetupSuite() {
	s.mongo = docker.GetMongo()
}

func (s *mongoSuite) SetupTest() {
}

func (s *mongoSuite) TearDownTest() {
}

func (s *mongoSuite) TearDownSuite() {
}

func (s *mongoSuite) TestMongo() {
	co := s.mongo.DB("practice").C("review")

	now := time.Now().Unix()
	type Review struct {
		ID         uuid.UUID
		WorkName   string
		Title      string
		Content    string
		CreateDate int64
		EditDate   int64
	}

	r1 := Review{
		ID:       uuid.New(),
		WorkName: "xd",
		Title:    "qq",
		Content:  "gg",
		CreaDate: now,
		EditDate: now,
	}
	c.Insert(r1)

	result := Review{}
	it := c.Find(bson.M{"ID": id.String()}).Limit(1).Iter()
	for it.Next(&result) {
		s.Equal(r1, result)
	}
}
