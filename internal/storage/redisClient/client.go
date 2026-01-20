package redisClient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/service-info-aggregator/internal/config"
)

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.Username,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
