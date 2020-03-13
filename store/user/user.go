package user

import (
	"encoding/json"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Create(context ctx.CTX, user *models.User) error
	GetByToken(context ctx.CTX, token string) (*models.User, error)
}

type impl struct {
	sql   *sqlx.DB
	redis redis.Service
}

func NewStore(db *sqlx.DB, redis redis.Service) Store {
	return &impl{
		sql:   db,
		redis: redis,
	}
}

func (u *impl) Create(context ctx.CTX, user *models.User) error {
	user.ID = uuid.New().String()
	user.Token = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("Create failed at bcrypt.GenerateFromPassword")
		return err
	}

	if _, err := u.sql.Exec("INSERT INTO users (id, token, email, password, register_date) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Token, user.Email, hashedPassword, time.Now().Unix()); err != nil {
		context.WithFields(logrus.Fields{
			"err":  err,
			"user": user,
		}).Error("Create failed at sql.Exec")
		return err
	}
	return nil
}

func (u *impl) GetByToken(context ctx.CTX, token string) (*models.User, error) {
	user := &models.User{}
	val, err := u.redis.Get(context, token)
	if err != nil {
		if err := u.sql.Get(user, "SELECT email, id FROM users WHERE token = $1", token); err != nil {
			context.WithField("err", err).Error("GetByToken failed at sql.Get")
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
