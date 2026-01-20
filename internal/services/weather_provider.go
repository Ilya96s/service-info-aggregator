package services

import (
	"context"

	"github.com/service-info-aggregator/internal/models/dto"
)

type WeatherProvider struct {
}

func (w *WeatherProvider) Name() string {
	return "weather"
}

func (w *WeatherProvider) CacheKey(c string) string {
	return "weather:" + c
}

func (w *WeatherProvider) Fetch(ctx context.Context, city string) (any, error) {
	return dto.WeatherResponse{
		City: city,
		Temp: 20,
	}, nil
}
