package user

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("Email address has been used")
)

type Store interface {
	Create(context ctx.CTX, user *models.User) (*models.User, error)
	GetByToken(context ctx.CTX, token uuid.UUID) (*models.User, error)
}

type impl struct {
	db    *sqlx.DB
	redis redis.Service
}

func New(db *sqlx.DB, redis redis.Service) Store {
	return &impl{
		db:    db,
		redis: redis,
	}
}

func (u *impl) Create(context ctx.CTX, user *models.User) (*models.User, error) {
	user.ID = uuid.New()
	user.Token = uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		context.WithField("err", err).Error("Create failed at bcrypt.GenerateFromPassword")
		return nil, err
	}

	query := "INSERT INTO users (id, token, name, email, password, register_date) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err := u.db.Exec(query, user.ID, user.Token, user.Name, user.Email, hashedPassword, time.Now().Unix()); err != nil {
		context.WithField("err", err).Error("Create failed at db.Exec")
		return nil, wrapError(err)
	}
	return user, nil
}

func (u *impl) GetByToken(context ctx.CTX, token uuid.UUID) (*models.User, error) {
	user := &models.User{}
	tokenString := token.String()
	if err := u.redis.Get(context, tokenString, user); err != nil {
		if err := u.db.Get(user, "SELECT email, id FROM users WHERE token = $1", token); err != nil {
			context.WithField("err", err).Error("GetByToken failed at db.Get")
			return nil, err
		}

		// Cache user in redis
		b, _ := json.Marshal(user)
		if err := u.redis.Set(context, tokenString, b, 7*24*time.Hour); err != nil {
			context.WithField("err", err).Error("GetByToken failed at redis.Set")
			// Fail but still acceptable result, so no return here
		}
	}
	return user, nil
}

func wrapError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		// Code 23505 represents unique_violation
		if pqErr.Code == "23505" {
			return ErrDuplicateEmail
		}
	}
	return err
}
