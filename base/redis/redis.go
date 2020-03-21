package redis

import (
	"errors"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
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
	if len(key) == 0 {
		return nil, errors.New("empty key")
	}
	conn := r.redis.Get()
	defer Close(context, "Get", conn)

	val, err := conn.Do("GET", key)
	if err != nil {
		context.WithFields(logrus.Fields{
			"err": err,
			"key": key,
		}).Error("Get failed at conn.Do")
		return nil, err
	}

	res, ok := val.([]byte)
	if !ok {
		return nil, errors.New("type assertion failed")
	}
	return res, nil
}

func (r *impl) Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	conn := r.redis.Get()
	defer Close(context, "Set", conn)

	_, err := conn.Do("SET", key, value, "EX", int64(expiration))
	if err != nil {
		context.WithFields(logrus.Fields{
			"err":        err,
			"key":        key,
			"value":      value,
			"expiration": expiration,
		}).Error("Set failed at conn.Do")
		return err
	}

	return nil
}
