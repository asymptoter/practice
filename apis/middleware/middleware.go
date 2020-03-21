package middleware

import (
	"net/http"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/store/user"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUser(us user.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := ctx.Background()
		token := c.Request.Header.Get("token")

		user, err := us.GetByToken(context, uuid.MustParse(token))
		if err != nil {
			context.WithFields(logrus.Fields{
				"err":   err,
				"token": token,
			}).Error("GetUser failed at GetByToken")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Set("userInfo", user)
		c.Next()
	}
}
