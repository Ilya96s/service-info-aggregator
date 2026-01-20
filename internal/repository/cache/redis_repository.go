package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: client,
	}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.redisClient.Set(ctx, key, value, ttl).Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}
