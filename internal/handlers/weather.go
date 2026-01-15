package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/service-info-aggregator/internal/models"
	"github.com/service-info-aggregator/internal/repository/cache"
	"github.com/service-info-aggregator/internal/services"
)

const KeyPrefix = "weather:"

type WeatherHandler struct {
	weatherService  *services.WeatherService
	redisRepository *cache.RedisRepository
	ttl             time.Duration
}

func NewWeatherHandler(ws *services.WeatherService, repo *cache.RedisRepository, ttl time.Duration) *WeatherHandler {
	return &WeatherHandler{
		weatherService:  ws,
		redisRepository: repo,
		ttl:             ttl,
	}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WeatherHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city parameter is required", http.StatusBadRequest)
		return
	}

	key := KeyPrefix + city

	data, err := h.redisRepository.Get(ctx, key)
	if err == nil {
		var resp models.WeatherResponse
		if err := json.Unmarshal([]byte(data), &resp); err != nil {
			slog.Error("failed to unmarshal weather from redis", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		responseWithJSON(w, http.StatusOK, resp)
		return
	} else if errors.Is(err, redis.Nil) {
		resp, err := h.weatherService.GetWeather(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(resp)
		if err == nil {
			h.redisRepository.Set(ctx, key, string(bytes), h.ttl)
		} else {
			slog.Error("failed to marshal weather for redis", "error", err)
		}

		responseWithJSON(w, http.StatusOK, resp)
		return
	} else {
		slog.Error("failed to get weather from redis", "key", key, "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func responseWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
