package cache

import (
	"context"
	"fmt"
	"time"

	"jarvis/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and returns a new Redis client.
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}
