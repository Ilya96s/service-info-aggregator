package services

import (
	"context"

	"github.com/service-info-aggregator/internal/models/dto"
	"github.com/service-info-aggregator/internal/repository/postgres"
)

type PopularDataService struct {
	Repo *postgres.PopularDataRepository
}

func NewPopularDataService(repo *postgres.PopularDataRepository) *PopularDataService {
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
