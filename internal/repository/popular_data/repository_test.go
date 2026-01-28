package popular_data_test

import (
	"context"
	"service-info-aggregator/internal/repository/popular_data"
	"testing"

	"service-info-aggregator/internal/model/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPopularDataRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := popular_data.NewPopularDataRepository(db)

	input := &dto.PopularDataDto{
		DataType: "weather",
		Key:      "Moscow",
	}

	rows := sqlmock.NewRows([]string{"id", "data_type", "key"}).AddRow(1, "weather", "Moscow")
	mock.ExpectQuery("INSERT INTO popular_data").
		WithArgs(input.DataType, input.Key, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)
	result, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "weather", result.DataType)
	assert.Equal(t, "Moscow", result.Key)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestPopularDataRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := popular_data.NewPopularDataRepository(db)
	rows := mock.NewRows([]string{"data_type", "key"}).AddRow("weather", "Moscow")
	mock.ExpectQuery("SELECT data_type, key FROM popular_data").WillReturnRows(rows)
	result, err := repo.GetAll(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPopularDataRepository_GetById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := popular_data.NewPopularDataRepository(db)
	rows := mock.NewRows([]string{"data_type", "key"}).AddRow("weather", "Novosibirsk")
	mock.ExpectQuery("SELECT data_type, key FROM popular_data").WillReturnRows(rows)
	result, err := repo.GetById(context.Background(), 0)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPopularDataRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := popular_data.NewPopularDataRepository(db)
	input := &dto.PopularDataDto{
		DataType: "weather",
		Key:      "Berlin",
	}
	rows := mock.NewRows([]string{"id", "data_type", "key"}).AddRow(1, "weather", "Berlin")
	mock.ExpectQuery("UPDATE popular_data SET").WillReturnRows(rows)
	result, err := repo.Update(context.Background(), 1, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "weather", result.DataType)
	assert.Equal(t, "Berlin", result.Key)
	require.NoError(t, mock.ExpectationsWereMet())
}
