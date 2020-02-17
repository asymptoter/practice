package middleware

import (
	"net/http"

	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/store/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUser(us user.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := ctx.Background()
		data := &struct {
			Token string `json:"token"`
		}{}
		if err := c.ShouldBind(data); err != nil {
			context.WithFields(logrus.Fields{
				"params": data.Token,
				"error":  err,
			}).Error("GetUser failed at ShouldBind ", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		user, err := us.GetByToken(context, data.Token)
		if err != nil {
			context.WithFields(logrus.Fields{
				"err":   err,
				"token": data.Token,
			}).Error("GetUser failed at GetByToken")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Set("userInfo", user)
		c.Next()
	}
}
