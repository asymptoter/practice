package mongo

import (
	"gopkg.in/mgo.v2"
)

func MustNew(connectionString string) *mgo.Session {
	res, err := New(connectionString)
	if err != nil {
		panic("New mongo failed")
	}
	return res
}

func New(connectionString string) (*mgo.Session, error) {
	dialInfo, err := mgo.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	return session, nil
}
