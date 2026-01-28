package popular_data

import (
	"context"

	"service-info-aggregator/internal/model/dto"
	"service-info-aggregator/internal/repository/popular_data"
)

type PopularDataService struct {
	Repo popular_data.Repository
}

func NewPopularDataService(repo popular_data.Repository) *PopularDataService {
	return &PopularDataService{
		Repo: repo,
	}
}

func (s *PopularDataService) Create(ctx context.Context, popularDataDto *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	return s.Repo.Create(ctx, popularDataDto)
}

func (s *PopularDataService) GetAll(ctx context.Context) ([]dto.PopularDataDto, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PopularDataService) GetById(ctx context.Context, id int) (*dto.PopularDataDto, error) {
	return s.Repo.GetById(ctx, id)
}

func (s *PopularDataService) Update(ctx context.Context, id int, inputData *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	return s.Repo.Update(ctx, id, inputData)
}

func (s *PopularDataService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
