package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379" // fallback for dev
	}

	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}
