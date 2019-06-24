package utils

import (
	"fmt"
)

type Context struct {
	context.Context
}

func GetContext(c *gin.Context) *Context {
	ctx := c.Request.Context()
}
