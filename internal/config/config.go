package config

import (
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
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

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:         getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		Username:     getEnv("REDIS_USERNAME", ""),
		Password:     getEnv("REDIS_PASSWORD", ""),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		DB:           getEnvInt("REDIS_DB", 0),
		DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 5*time.Second),
		WeatherTTL:   getEnvDuration("REDIS_WEATHER_TTL", 3000*time.Second),
	}
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9091")},
		Topic:   getEnv("KAFKA_TOPIC", "external.events"),
		GroupID: getEnv("KAFKA_GROUP_ID", "aggregator"),
	}
}

type PostgresConfig struct {
	Host            string
	Port            int
	DBName          string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:            getEnv("POSTGRES_HOST", "localhost"),
		Port:            getEnvInt("POSTGRES_PORT", 5432),
		DBName:          getEnv("POSTGRES_DB", "service_info"),
		User:            getEnv("POSTGRES_USER", "postgres"),
		Password:        getEnv("POSTGRES_PASSWORD", "postgres"),
		SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
		MaxOpenConns:    getEnvInt("POSTGRES_MAX_OPEN_CONNS", 10),
		MaxIdleConns:    getEnvInt("POSTGRES_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDuration("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute),
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
