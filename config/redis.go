package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Ctx
var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	// Connect to Redis
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
