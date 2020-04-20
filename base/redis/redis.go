package redis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
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

func Close(context ctx.CTX, functionName string, conn redis.Conn) {
	if err := conn.Close(); err != nil {
		context.WithField("err", err).Error(functionName + " failed at Close")
	}
}

func (r *impl) Get(context ctx.CTX, key string, res interface{}) error {
	context = ctx.WithValue(context, "key", key)
	if len(key) == 0 {
		return errors.New("empty key")
	}
	conn := r.redis.Get()
	defer Close(context, "Get", conn)

	val, err := conn.Do("GET", key)
	if err != nil {
		context.WithField("err", err).Error("Get failed at conn.Do")
		return err
	}

	if err := json.Unmarshal(val.([]byte), res); err != nil {
		context.WithField("err", err).Error("Get failed at json.Unmarshal")
		return err
	}
	return nil
}

func (r *impl) Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	context = ctx.WithValue(context, "key", key)
	conn := r.redis.Get()
	defer Close(context, "Set", conn)

	v, err := json.Marshal(value)
	if err != nil {
		context.WithField("err", err).Error("Set failed at json.Masharl")
		return err
	}

	if _, err := conn.Do("SET", key, v, "EX", int64(expiration)); err != nil {
		context.WithFields(logrus.Fields{
			"err":        err,
			"expiration": expiration,
		}).Error("Set failed at conn.Do")
		return err
	}

	return nil
}
