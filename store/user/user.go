package user

import (
	"encoding/json"
	"time"

	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/asymptoter/geochallenge-backend/base/redis"
	"github.com/asymptoter/geochallenge-backend/models"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	GetByToken(context ctx.CTX, token string) (*models.User, error)
}

type impl struct {
	mysql *sqlx.DB
	redis redis.Service
}

func NewStore(db *sqlx.DB, redis redis.Service) Store {
	return &impl{
		mysql: db,
		redis: redis,
	}
}

func (u *impl) GetByToken(context ctx.CTX, token string) (*models.User, error) {
	user := &models.User{}
	val, err := u.redis.Get(context, token)
	if err != nil {
		if err := u.mysql.Get(user, "SELECT email, id, activation_number from users where token = ?", token); err != nil {
			context.WithField("err", err).Error("GetByToken failed at mysql.Get")
			return nil, err
		}

		// Cache user in redis
		b, _ := json.Marshal(user)
		if err := u.redis.Set(context, token, b, 7*24*time.Hour); err != nil {
			context.WithField("err", err).Error("GetByToken failed at redis.Set")
			// Fail but still acceptable result, so no return here
		}
	} else if err := json.Unmarshal([]byte(val), user); err != nil {
		context.WithField("err", err).Error("GetUserByToken failed at Unmarshal")
		return nil, err
	}

	return user, nil
}
