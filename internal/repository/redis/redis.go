package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"template/internal/config"
)

func NewRedisDB(cfg *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	if res := client.Ping(context.TODO()); res.Err() != nil {
		if err := client.Close(); err != nil {
			return nil, fmt.Errorf("failed to disconnect Redis after ping failure: %w", err)
		}
		return nil, fmt.Errorf("redis ping failed: %w", res.Err())
	}

	return client, nil
}
