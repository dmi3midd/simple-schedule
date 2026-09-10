package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoWeek = errors.New("no week found")
)

type WeekRepository interface {
	// GetById retrieves a week by its ID.
	// Returns [ErrNoWeek] if the week is not found.
	GetById(ctx context.Context, id uuid.UUID) (*domain.Week, error)
	// GetAll retrieves all weeks
	// Returns an empty slice if no weeks are found.
	GetAll(ctx context.Context) ([]domain.Week, error)
	// Create creates a new week
	Create(ctx context.Context, week domain.Week) (*domain.Week, error)
	// Update updates an existing week
	Update(ctx context.Context, week domain.Week) (*domain.Week, error)
	// Delete deletes a week
	Delete(ctx context.Context, id uuid.UUID) error
}

type weekRepository struct {
	db *sqlx.DB
}

func NewWeekRepository(db *sqlx.DB) WeekRepository {
	return &weekRepository{db: db}
}

func (r *weekRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Week, error) {
	op := "WeekRepository.GetById"
	query := `
		SELECT id, title, created_at, updated_at 
		FROM weeks 
		WHERE id = $1`
	var week domain.Week
	executor := ExtractTx(ctx, r.db)
	err := sqlx.GetContext(ctx, executor, &week, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoWeek)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &week, nil
}

func (r *weekRepository) GetAll(ctx context.Context) ([]domain.Week, error) {
	op := "WeekRepository.GetAll"
	query := `
		SELECT id, title, created_at, updated_at 
		FROM weeks`
	var weeks []domain.Week
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &weeks, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return weeks, nil
}

func (r *weekRepository) Create(ctx context.Context, week domain.Week) (*domain.Week, error) {
	op := "WeekRepository.Create"
	query := `
	INSERT INTO weeks (id, title, created_at, updated_at) 
	VALUES (:id, :title, :created_at, :updated_at) 
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, week)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &week, nil
}

func (r *weekRepository) Update(ctx context.Context, week domain.Week) (*domain.Week, error) {
	op := "WeekRepository.Update"
	query := `
	UPDATE weeks 
	SET title = :title, updated_at = :updated_at
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, week)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &week, nil
}

func (r *weekRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "WeekRepository.Delete"
	query := `
	DELETE FROM weeks 
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
