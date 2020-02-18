package redis

import (
	"time"

	"github.com/asymptoter/geochallenge-backend/base/config"
	"github.com/asymptoter/geochallenge-backend/base/ctx"
	"github.com/go-redis/redis/v7"
)

type Service interface {
	Get(context ctx.CTX, key string) ([]byte, error)
	Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error
}

type impl struct {
	redis *redis.Client
}

func NewService() Service {
	cfg := config.Value.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &impl{
		redis: client,
	}
}

func (r *impl) Get(context ctx.CTX, key string) ([]byte, error) {
	return []byte{}, nil
}

func (r *impl) Set(context ctx.CTX, key string, value interface{}, expiration time.Duration) error {
	//redis.Set("a", "b", time.Hour)
	return nil
}
