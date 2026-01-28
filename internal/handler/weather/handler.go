package weather

import (
	"encoding/json"
	"net/http"
	"time"

	"service-info-aggregator/internal/repository/aggregation_data"
	"service-info-aggregator/internal/service/aggregation"
)

type WeatherHandler struct {
	aggregationService *aggregation.AggregationService
	weatherProvider    *aggregation.WeatherProvider
	cache              *aggregation_data.RedisRepository
	ttl                time.Duration
}

func NewWeatherHandler(aggregationService *aggregation.AggregationService, weatherProvider *aggregation.WeatherProvider,
	repo *aggregation_data.RedisRepository, ttl time.Duration) *WeatherHandler {
	return &WeatherHandler{
		aggregationService: aggregationService,
		weatherProvider:    weatherProvider,
		cache:              repo,
		ttl:                ttl,
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

	key := h.weatherProvider.CacheKey(city)

	if cached, err := h.cache.Get(ctx, key); err == nil {
		responseWithJSON(w, http.StatusOK, cached)
		return
	}

	data, err := h.aggregationService.Execute(ctx, h.weatherProvider, city)
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
