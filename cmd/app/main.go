package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"service-info-aggregator/internal/background"
	"service-info-aggregator/internal/config"
	"service-info-aggregator/internal/handler/popular_data"
	"service-info-aggregator/internal/handler/weather"
	"service-info-aggregator/internal/messaging/kafka"
	"service-info-aggregator/internal/repository/aggregation_data"
	postgresRepo "service-info-aggregator/internal/repository/popular_data"
	"service-info-aggregator/internal/service/aggregation"
	popular_data2 "service-info-aggregator/internal/service/popular_data"
	"service-info-aggregator/internal/storage/postgres"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// --- Конфиги ---
	pgCfg := config.NewPostgresConfig()
	redisCfg := config.NewRedisConfig()
	kafkaCfg := config.NewKafkaConfig()

	// --- Postgres ---
	db, err := postgres.NewPostgresConnection(pgCfg)
	if err != nil {
		slog.Error("failed to connect postgres", "error", err)
		return
	}

	// --- Redis ---
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Username: redisCfg.Username,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return
	}
	repo := aggregation_data.NewRedisRepository(rdb)

	// --- Kafka Producer ---
	producer, err := kafka.NewKafkaProducer(
		[]string{"127.0.0.1:9091", "127.0.0.1:9092", "127.0.0.1:9093"},
		"aggregator-producer",
	)
	if err != nil {
		slog.Error("failed to create kafka producer", "error", err)
		return
	}
	defer producer.Close()

	// --- Popular Data Repository ---
	popularDataRepository := postgresRepo.NewPopularDataRepository(db)

	// --- Popular Data Service
	popularDataService := popular_data2.NewPopularDataService(popularDataRepository)

	// --- Сервис агрегирования ---
	aggService := aggregation.NewAggregationService(producer, kafkaCfg.Topic)

	// --- Popular Data Handler ---
	popularDataHandler := popular_data.NewPopularDataHandler(popularDataService)

	// --- Weather Provider ---
	weatherProvider := &aggregation.WeatherProvider{}

	// --- Weather Handler (HTTP) ---
	weatherHandler := weather.NewWeatherHandler(aggService, weatherProvider, repo, redisCfg.WeatherTTL)

	mux := http.NewServeMux()

	mux.Handle("/weather", weatherHandler)
	mux.HandleFunc("/popular-data", popularDataHandler.HandleCollection)
	mux.HandleFunc("/popular-data/", popularDataHandler.HandleItem)

	// --- Kafka Event Handler ---
	weatherEventHandler := kafka.NewWeatherEventHandler(repo, redisCfg.WeatherTTL)
	eventRouter := kafka.NewEventRouter(weatherEventHandler)

	// --- Kafka Consumer ---
	consumer, err := kafka.NewKafkaConsumer(
		[]string{"127.0.0.1:9091", "127.0.0.1:9092", "127.0.0.1:9093"},
		"aggregator-consumer",
		eventRouter,
	)
	if err != nil {
		slog.Error("failed to create kafka consumer", "error", err)
		return
	}
	defer consumer.Close()

	// --- Scheduler ---
	scheduler := background.NewPriorityScheduler(popularDataService, aggService, 30*time.Second)

	go func() {
		scheduler.Start(ctx)
	}()

	// --- Запуск Kafka Consumer в отдельной горутине ---
	go func() {
		if err := consumer.Run(ctx, []string{kafkaCfg.Topic}); err != nil {
			slog.Error("kafka consumer stopped", "error", err)
		}
	}()

	// --- Запуск HTTP сервера ---
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("HTTP server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server stopped", "error", err)
			cancel()
		}
	}()

	// --- Ждём сигнала выхода ---
	<-ctx.Done()
	slog.Info("shutting down...")

	// --- Грейсфул остановка HTTP сервера ---
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}
