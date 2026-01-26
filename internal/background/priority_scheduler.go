package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/service-info-aggregator/internal/service"
)

type PriorityScheduler struct {
	popularDataService *service.PopularDataService
	aggregationService *service.AggregationService
	interval           time.Duration
}

func NewPriorityScheduler(ps *service.PopularDataService, as *service.AggregationService, interval time.Duration) *PriorityScheduler {
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
			_, err := s.aggregationService.Execute(ctx, &service.WeatherProvider{}, item.Key)
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
