package popular_data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"service-info-aggregator/internal/model/dto"
)

type Repository interface {
	Create(ctx context.Context, dto *dto.PopularDataDto) (*dto.PopularDataDto, error)
	GetAll(ctx context.Context) ([]dto.PopularDataDto, error)
	GetById(ctx context.Context, id int) (*dto.PopularDataDto, error)
	Update(ctx context.Context, id int, dto *dto.PopularDataDto) (*dto.PopularDataDto, error)
	Delete(ctx context.Context, id int) error
}

type PopularDataRepository struct {
	db *sql.DB
}

func NewPopularDataRepository(db *sql.DB) *PopularDataRepository {
	return &PopularDataRepository{
		db: db,
	}
}

func (r *PopularDataRepository) Create(ctx context.Context, inputData *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	query := `
		INSERT INTO popular_data (data_type, key, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, data_type, key
	`

	var created dto.PopularDataDto
	err := r.db.QueryRowContext(ctx, query,
		inputData.DataType,
		inputData.Key,
		time.Now(),
		time.Now(),
	).Scan(
		&created.ID,
		&created.DataType,
		&created.Key,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *PopularDataRepository) GetAll(ctx context.Context) ([]dto.PopularDataDto, error) {
	query := `SELECT data_type, key FROM popular_data`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.PopularDataDto, 0)
	for rows.Next() {
		var popularDataDto dto.PopularDataDto
		if err := rows.Scan(&popularDataDto.DataType, &popularDataDto.Key); err != nil {
			return nil, err
		}
		results = append(results, popularDataDto)
	}

	return results, rows.Err()
}

func (r *PopularDataRepository) GetById(ctx context.Context, id int) (*dto.PopularDataDto, error) {
	query := `SELECT data_type, key FROM popular_data WHERE id = $1`

	var result dto.PopularDataDto

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&result.DataType, &result.Key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &result, err
}

func (r *PopularDataRepository) Update(ctx context.Context, id int, inputData *dto.PopularDataDto) (*dto.PopularDataDto, error) {
	query := `UPDATE popular_data 
			  SET data_type = $1, key = $2, updated_at = $3 
			  WHERE id = $4
			  RETURNING id, data_type, key`

	var updated dto.PopularDataDto
	err := r.db.QueryRowContext(ctx, query, inputData.DataType, inputData.Key, time.Now(), id).Scan(&updated.ID, &updated.DataType, &updated.Key)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *PopularDataRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM popular_data WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
