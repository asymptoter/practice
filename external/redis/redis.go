package redis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"

	"github.com/gomodule/redigo/redis"
)

type Service interface {
	Get(context ctx.CTX, key string, res interface{}) error
	Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error
}

type impl struct {
	redis *redis.Pool
}

func NewService(address string) Service {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", address) },
	}

	return &impl{
		redis: pool,
	}
}

func Close(ctx ctx.CTX, functionName string, conn redis.Conn) {
	if err := conn.Close(); err != nil {
		ctx.With("err", err).Error(functionName + " failed at Close")
	}
}

func (r *impl) Get(ctx ctx.CTX, key string, res interface{}) error {
	ctx = ctx.With("key", key)
	if len(key) == 0 {
		return errors.New("empty key")
	}
	conn := r.redis.Get()
	defer Close(ctx, "Get", conn)

	val, err := conn.Do("GET", key)
	if err != nil {
		ctx.Error(err)
		return err
	}

	if err := json.Unmarshal(val.([]byte), res); err != nil {
		ctx.Error(err)
		return err
	}
	return nil
}

func (r *impl) Set(ctx ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	ctx = ctx.With("key", key)
	conn := r.redis.Get()
	defer Close(ctx, "Set", conn)

	v, err := json.Marshal(value)
	if err != nil {
		ctx.Error(err)
		return err
	}

	if _, err := conn.Do("SET", key, v, "EX", int64(expiration)); err != nil {
		ctx.With("expiration", expiration).Error(err)
		return err
	}

	return nil
}
