package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/service-info-aggregator/internal/repository/cache"
)

type WeatherEventHandler struct {
	cache *cache.RedisRepository
	ttl   time.Duration
}

func NewWeatherEventHandler(c *cache.RedisRepository, ttl time.Duration) *WeatherEventHandler {
	return &WeatherEventHandler{
		cache: c,
		ttl:   ttl,
	}
}

func (h *WeatherEventHandler) Type() string {
	return "weather"
}

func (h *WeatherEventHandler) Handle(ctx context.Context, key string, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal weather payload: %w", err)
	}

	cacheKey := "weather:" + key
	err = h.cache.Set(ctx, cacheKey, string(bytes), h.ttl)
	if err != nil {
		slog.Error("Redis Set failed", "error", err)
		return err
	}
	return nil
}
