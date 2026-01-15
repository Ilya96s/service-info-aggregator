package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr         string
	Username     string
	Password     string
	DB           int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	WeatherTTL   time.Duration
}

func NewConfig() *Config {
	return &Config{
		Addr:         getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Username:     getEnv("REDIS_USERNAME", ""),
		Password:     getEnv("REDIS_PASSWORD", ""),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		DB:           getEnvInt("REDIS_DB", 0),
		DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 5*time.Second),
		WeatherTTL:   getEnvDuration("REDIS_WEATHER_TTL", 30*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
