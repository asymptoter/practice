package ctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CTX struct {
	context.Context
	Fields logrus.Fields
}

func NewContext(c *gin.Context) CTX {
	return c.Request.Context().(CTX)
}

func (c *CTX) WithValue(key string, val interface{}) {
	c.Fields[key] = val
}

func (c *CTX) WithField(key string, val interface{}) *logrus.Entry {
	return logrus.WithFields(c.Fields).WithField(key, val)
}

func (c *CTX) WithFields(f logrus.Fields) *logrus.Entry {
	return logrus.WithFields(c.Fields).WithFields(f)
}
