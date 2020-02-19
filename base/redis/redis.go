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

func Close(context ctx.CTX, functionName string, conn redis.Conn) {
	if err := conn.Close(); err != nil {
		context.WithField("err", err).Error(functionName + " failed at Close")
	}
}

func (r *impl) Get(context ctx.CTX, key string) ([]byte, error) {
	conn := r.redis.Get()
	val, err := conn.Do("GET", key)
	defer Close(context, "Get", conn)
	return val.([]byte), err
}

func (r *impl) Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	conn := r.redis.Get()
	_, err := conn.Do("SET", key, value, expiration)
	defer Close(context, "Set", conn)
	return err
}
