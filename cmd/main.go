package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/service-info-aggregator/internal/config"
	"github.com/service-info-aggregator/internal/handlers"
	"github.com/service-info-aggregator/internal/repository/cache"
	"github.com/service-info-aggregator/internal/services"
	"github.com/service-info-aggregator/internal/storage/rediss"
)

func main() {
	// Logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Config
	cfg := config.NewConfig()

	ctx := context.Background()

	// Redis
	client, err := rediss.NewRedisClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to init redis", "error", err)
		panic(err)
	}

	repo := cache.NewRedisRepository(client)

	// Weather service & handler
	weatherService := services.NewWeatherService(repo)
	weatherHandler := handlers.NewWeatherHandler(weatherService, repo, cfg.WeatherTTL)

	mux := http.NewServeMux()
	mux.Handle("/weather/", weatherHandler)

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
