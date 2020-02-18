package redis

import (
	"time"

	"github.com/asymptoter/geochallenge-backend/base/config"
	"github.com/asymptoter/geochallenge-backend/base/ctx"

	"github.com/gomodule/redigo/redis"
)

type Service interface {
	Get(context ctx.CTX, key string) ([]byte, error)
	Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error
}

type impl struct {
	redis *redis.Pool
}

func NewService() Service {
	cfg := config.Value.Redis

	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", cfg.Address) },
	}

	return &impl{
		redis: pool,
	}
}

func (r *impl) Get(context ctx.CTX, key string) ([]byte, error) {
	val, err := r.redis.Get().Do("GET", key)
	return val.([]byte), err
}

func (r *impl) Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	_, err := r.redis.Get().Do("SET", key, value, expiration)
	return err
}
