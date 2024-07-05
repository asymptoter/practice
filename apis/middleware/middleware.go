package middleware

import (
	"net/http"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/store/user"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func GetUser(us user.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctx.Background()
		token := c.Request.Header.Get("token")

		user, err := us.GetByToken(ctx, uuid.MustParse(token))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Set("userInfo", user)
		c.Next()
	}
}
