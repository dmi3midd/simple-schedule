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
	ErrNoSlot = errors.New("no slot found")
)

type SlotRepository interface {
	// GetById retrieves a slot by its ID.
	// Returns [ErrNoSlot] if the slot is not found.
	GetById(ctx context.Context, id uuid.UUID) (*domain.Slot, error)
	// GetAll retrieves all slots.
	// Returns an empty slice if no slots are found.
	GetAll(ctx context.Context) ([]domain.Slot, error)
	// GetByWeekId retrieves all slots for a given week.
	// Returns an empty slice if no slots are found.
	GetByWeekId(ctx context.Context, weekId uuid.UUID) ([]domain.Slot, error)
	// GetByWeekIdAndDayOfWeek retrieves all slots for a given week and day of week.
	// Returns an empty slice if no slots are found.
	GetByWeekIdAndDayOfWeek(ctx context.Context, weekId uuid.UUID, dayOfWeek domain.DayOfWeek) ([]domain.Slot, error)
	// Create creates a new slot
	Create(ctx context.Context, slot domain.Slot) (*domain.Slot, error)
	// Update updates an existing slot
	Update(ctx context.Context, slot domain.Slot) (*domain.Slot, error)
	// Delete deletes a slot
	Delete(ctx context.Context, id uuid.UUID) error
}

type slotRepository struct {
	db *sqlx.DB
}

func NewSlotRepository(db *sqlx.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Slot, error) {
	op := "SlotRepository.GetById"
	query := `
		SELECT id, week_id, day_of_week, activity, start_time, end_time, created_at, updated_at 
		FROM slots 
		WHERE id = $1
	`
	var slot domain.Slot
	executor := ExtractTx(ctx, r.db)
	err := sqlx.GetContext(ctx, executor, &slot, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoSlot)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &slot, nil
}

func (r *slotRepository) GetAll(ctx context.Context) ([]domain.Slot, error) {
	op := "SlotRepository.GetAll"
	query := `
		SELECT id, week_id, day_of_week, activity, start_time, end_time, created_at, updated_at 
		FROM slots
		ORDER BY day_of_week, start_time
	`
	var slots []domain.Slot
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &slots, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return slots, nil
}

func (r *slotRepository) GetByWeekId(ctx context.Context, weekId uuid.UUID) ([]domain.Slot, error) {
	op := "SlotRepository.GetByWeekId"
	query := `
		SELECT id, week_id, day_of_week, activity, start_time, end_time, created_at, updated_at 
		FROM slots
		WHERE week_id = $1
		ORDER BY day_of_week, start_time
	`
	var slots []domain.Slot
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &slots, query, weekId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return slots, nil
}

func (r *slotRepository) GetByWeekIdAndDayOfWeek(ctx context.Context, weekId uuid.UUID, dayOfWeek domain.DayOfWeek) ([]domain.Slot, error) {
	op := "SlotRepository.GetByWeekIdAndDayOfWeek"
	query := `
		SELECT id, week_id, day_of_week, activity, start_time, end_time, created_at, updated_at 
		FROM slots
		WHERE week_id = $1 AND day_of_week = $2
		ORDER BY start_time
	`
	var slots []domain.Slot
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &slots, query, weekId, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return slots, nil
}

func (r *slotRepository) Create(ctx context.Context, slot domain.Slot) (*domain.Slot, error) {
	op := "SlotRepository.Create"
	query := `
		INSERT INTO slots (id, week_id, day_of_week, activity, start_time, end_time, created_at, updated_at) 
		VALUES (:id, :week_id, :day_of_week, :activity, :start_time, :end_time, :created_at, :updated_at)
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, slot)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &slot, nil
}

func (r *slotRepository) Update(ctx context.Context, slot domain.Slot) (*domain.Slot, error) {
	op := "SlotRepository.Update"
	query := `
		UPDATE slots 
		SET week_id = :week_id, day_of_week = :day_of_week, activity = :activity, start_time = :start_time, end_time = :end_time, updated_at = :updated_at
		WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, slot)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &slot, nil
}

func (r *slotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "SlotRepository.Delete"
	query := `
		DELETE FROM slots 
		WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
