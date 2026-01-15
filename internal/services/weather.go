package services

import (
	"github.com/service-info-aggregator/internal/models"
	"github.com/service-info-aggregator/internal/repository/cache"
)

type WeatherService struct {
	redisRepository cache.RedisRepository
}

func NewWeatherService(redisRepository *cache.RedisRepository) *WeatherService {
	return &WeatherService{
		redisRepository: *redisRepository,
	}
}

// Заглушка
func (s *WeatherService) GetWeather(city string) (models.WeatherResponse, error) {
	return models.WeatherResponse{
		City:        city,
		Temperature: 20,
	}, nil
}
