package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/service-info-aggregator/internal/repository/cache"
	"github.com/service-info-aggregator/internal/services"
)

type WeatherHandler struct {
	service  *services.AggregationService
	provider *services.WeatherProvider
	cache    *cache.RedisRepository
	ttl      time.Duration
}

func NewWeatherHandler(service *services.AggregationService, provider *services.WeatherProvider,
	repo *cache.RedisRepository, ttl time.Duration) *WeatherHandler {
	return &WeatherHandler{
		service:  service,
		provider: provider,
		cache:    repo,
		ttl:      ttl,
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

	key := h.provider.CacheKey(city)

	if cached, err := h.cache.Get(ctx, key); err == nil {
		responseWithJSON(w, http.StatusOK, cached)
		return
	}

	data, err := h.service.Execute(ctx, h.provider, city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWithJSON(w, http.StatusOK, data)
}

func responseWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func responseWithError(w http.ResponseWriter, statusCode int, message string) {
	responseWithJSON(w, statusCode, map[string]string{"error": message})
}
