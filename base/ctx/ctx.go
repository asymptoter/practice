package ctx

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CTX struct {
	context.Context
	logrus.FieldLogger
}

func Background() CTX {
	return CTX{
		Context:     context.Background(),
		FieldLogger: logrus.StandardLogger(),
	}
}

func WithValue(parent CTX, key string, val interface{}) CTX {
	return CTX{
		Context:     context.WithValue(parent, key, val),
		FieldLogger: parent.FieldLogger.WithField(key, val),
	}
}

func WithValues(parent CTX, m map[string]interface{}) CTX {
	c := parent
	for k, v := range m {
		c = WithValue(c, k, v)
	}
	return c
}
