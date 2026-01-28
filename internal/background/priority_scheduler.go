package background

import (
	"context"
	"log/slog"
	"time"

	"service-info-aggregator/internal/service/aggregation"
	"service-info-aggregator/internal/service/popular_data"
)

type PriorityScheduler struct {
	popularDataService *popular_data.PopularDataService
	aggregationService *aggregation.AggregationService
	interval           time.Duration
}

func NewPriorityScheduler(ps *popular_data.PopularDataService, as *aggregation.AggregationService, interval time.Duration) *PriorityScheduler {
	return &PriorityScheduler{
		popularDataService: ps,
		aggregationService: as,
		interval:           interval,
	}
}

func (s *PriorityScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	slog.Info("priority scheduler started", "interval", s.interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("priority scheduler stopped")
			return
		case <-ticker.C:
			s.execute(ctx)
		}
	}
}

func (s *PriorityScheduler) execute(ctx context.Context) {
	items, err := s.popularDataService.GetAll(ctx)
	if err != nil {
		slog.Error("failed to fetch popular data items", "err", err)
		return
	}

	for _, item := range items {
		switch item.DataType {
		case "weather":
			_, err := s.aggregationService.Execute(ctx, &aggregation.WeatherProvider{}, item.Key)
			if err != nil {
				slog.Error("aggregation failed",
					"type", item.DataType,
					"key", item.Key,
					"error", err)
			}
		default:
			slog.Warn("unknown data type", "type", item.DataType)
		}
	}
}
