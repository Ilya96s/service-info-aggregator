package popular_data_test

import (
	"context"
	"errors"
	"service-info-aggregator/internal/model/dto"
	"service-info-aggregator/internal/service/popular_data"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPopularDataRepository struct {
	mock.Mock
}

func (m *MockPopularDataRepository) Create(ctx context.Context, d *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	args := m.Called(ctx, d)
	return args.Get(0).(*dto.PopularDataDto), args.Error(1)
}

func (m *MockPopularDataRepository) GetAll(ctx context.Context) ([]dto.PopularDataDto, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.PopularDataDto), args.Error(1)
}

func (m *MockPopularDataRepository) GetById(ctx context.Context, id int) (*dto.PopularDataDto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.PopularDataDto), args.Error(1)
}

func (m *MockPopularDataRepository) Update(ctx context.Context, id int, d *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	args := m.Called(ctx, id, d)
	return args.Get(0).(*dto.PopularDataDto), args.Error(1)
}

func (m *MockPopularDataRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPopularDataService_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockPopularDataRepository)
	service := popular_data.NewPopularDataService(repo)

	input := &dto.PopularDataDto{DataType: "weather", Key: "Moscow"}
	expected := &dto.PopularDataDto{ID: 1, DataType: "weather", Key: "Moscow"}

	repo.On("Create", ctx, input).Return(expected, nil)
	result, err := service.Create(ctx, input)

	require.NoError(t, err)
	require.Equal(t, expected, result)

	repo.AssertExpectations(t)
}

func TestPopularDataService_Create_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockPopularDataRepository)
	service := popular_data.NewPopularDataService(repo)

	input := &dto.PopularDataDto{}
	expectedErr := errors.New("db error")

	repo.On("Create", ctx, input).Return((*dto.PopularDataDto)(nil), expectedErr)

	result, err := service.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func TestPopularDataService_GetById_Success(t *testing.T) {
	ctx := context.Background()
	repo := new(MockPopularDataRepository)
	service := popular_data.NewPopularDataService(repo)

	expected := &dto.PopularDataDto{ID: 1, DataType: "weather", Key: "Moscow"}

	repo.On("GetById", ctx, 1).Return(expected, nil)

	result, err := service.GetById(ctx, 1)

	require.NoError(t, err)
	require.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
