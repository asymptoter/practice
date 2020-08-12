package ctx

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CTX struct {
	context.Context
	logrus.Entry
}

func Background() CTX {
	log := logrus.New()
	log.SetReportCaller(true)
	return CTX{
		Context: context.Background(),
		Entry:   *logrus.NewEntry(log),
	}
}

func WithValue(parent CTX, key string, val interface{}) CTX {
	return CTX{
		Context: context.WithValue(parent, key, val),
		Entry:   *parent.Entry.WithField(key, val),
	}
}

func WithValues(parent CTX, m logrus.Fields) CTX {
	return CTX{
		Context: parent.Context,
		Entry:   *parent.Entry.WithFields(m),
	}
}
