package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoSlot = errors.New("no slot found")
)

type SlotRepository interface {
	// GetById retrieves a slot by its ID with its associated tag.
	// Returns [ErrNoSlot] if the slot is not found.
	GetById(ctx context.Context, id uuid.UUID) (*domain.SlotWithTag, error)
	// GetByWeekId retrieves all slots for a given week with their associated tags.
	GetByWeekId(ctx context.Context, weekId uuid.UUID) ([]domain.SlotWithTag, error)
	// HasOverlap checks if there is any overlapping slot in the week on the same day.
	HasOverlap(ctx context.Context, weekId uuid.UUID, dayOfWeek domain.DayOfWeek, startTime, endTime int, excludeSlotId *uuid.UUID) (bool, error)
	// Create creates a new slot.
	Create(ctx context.Context, slot domain.Slot) error
	// Update updates an existing slot.
	Update(ctx context.Context, slot domain.Slot) error
	// Delete deletes a slot by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

type slotWithTagRow struct {
	ID          uuid.UUID        `db:"id"`
	WeekID      uuid.UUID        `db:"week_id"`
	TagID       *uuid.UUID       `db:"tag_id"`
	DayOfWeek   domain.DayOfWeek `db:"day_of_week"`
	Activity    string           `db:"activity"`
	StartTime   int              `db:"start_time"`
	EndTime     int              `db:"end_time"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
	TagRefID    *uuid.UUID       `db:"tag_ref_id"`
	TagTitle    *string          `db:"tag_title"`
	TagHexColor *string          `db:"tag_hex_color"`
	TagCreated  *time.Time       `db:"tag_created_at"`
	TagUpdated  *time.Time       `db:"tag_updated_at"`
}

func (r slotWithTagRow) toDomain() domain.SlotWithTag {
	slot := domain.Slot{
		ID:        r.ID,
		WeekID:    r.WeekID,
		TagID:     r.TagID,
		DayOfWeek: r.DayOfWeek,
		Activity:  r.Activity,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	var tag *domain.Tag
	if r.TagRefID != nil {
		title := ""
		if r.TagTitle != nil {
			title = *r.TagTitle
		}
		hexColor := ""
		if r.TagHexColor != nil {
			hexColor = *r.TagHexColor
		}
		createdAt := time.Time{}
		if r.TagCreated != nil {
			createdAt = *r.TagCreated
		}
		updatedAt := time.Time{}
		if r.TagUpdated != nil {
			updatedAt = *r.TagUpdated
		}
		tag = &domain.Tag{
			ID:        *r.TagRefID,
			Title:     title,
			HexColor:  hexColor,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return domain.SlotWithTag{
		Slot: slot,
		Tag:  tag,
	}
}

type slotRepository struct {
	db *sqlx.DB
}

func NewSlotRepository(db *sqlx.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.SlotWithTag, error) {
	op := "SlotRepository.GetById"
	query := `
		SELECT 
			s.id, s.week_id, s.tag_id, s.day_of_week, s.activity, s.start_time, s.end_time, s.created_at, s.updated_at,
			t.id AS tag_ref_id,
			t.title AS tag_title,
			t.hex_color AS tag_hex_color,
			t.created_at AS tag_created_at,
			t.updated_at AS tag_updated_at
		FROM slots s
		LEFT JOIN tags t ON s.tag_id = t.id
		WHERE s.id = $1
	`
	var row slotWithTagRow
	executor := ExtractTx(ctx, r.db)
	err := sqlx.GetContext(ctx, executor, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoSlot)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	res := row.toDomain()
	return &res, nil
}

func (r *slotRepository) GetByWeekId(ctx context.Context, weekId uuid.UUID) ([]domain.SlotWithTag, error) {
	op := "SlotRepository.GetByWeekId"
	query := `
		SELECT 
			s.id, s.week_id, s.tag_id, s.day_of_week, s.activity, s.start_time, s.end_time, s.created_at, s.updated_at,
			t.id AS tag_ref_id,
			t.title AS tag_title,
			t.hex_color AS tag_hex_color,
			t.created_at AS tag_created_at,
			t.updated_at AS tag_updated_at
		FROM slots s
		LEFT JOIN tags t ON s.tag_id = t.id
		WHERE s.week_id = $1
		ORDER BY s.day_of_week, s.start_time
	`
	var rows []slotWithTagRow
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &rows, query, weekId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]domain.SlotWithTag, len(rows))
	for i, row := range rows {
		result[i] = row.toDomain()
	}
	return result, nil
}

func (r *slotRepository) HasOverlap(ctx context.Context, weekId uuid.UUID, dayOfWeek domain.DayOfWeek, startTime, endTime int, excludeSlotId *uuid.UUID) (bool, error) {
	op := "SlotRepository.HasOverlap"
	query := `
		SELECT EXISTS(
			SELECT 1 FROM slots
			WHERE week_id = $1 
			  AND day_of_week = $2
			  AND start_time < $4 
			  AND end_time > $3
			  AND ($5::uuid IS NULL OR id != $5)
		)
	`
	var overlap bool
	executor := ExtractTx(ctx, r.db)
	err := sqlx.GetContext(ctx, executor, &overlap, query, weekId, dayOfWeek, startTime, endTime, excludeSlotId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return overlap, nil
}

func (r *slotRepository) Create(ctx context.Context, slot domain.Slot) error {
	op := "SlotRepository.Create"
	query := `
		INSERT INTO slots (id, week_id, tag_id, day_of_week, activity, start_time, end_time, created_at, updated_at) 
		VALUES (:id, :week_id, :tag_id, :day_of_week, :activity, :start_time, :end_time, :created_at, :updated_at)
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, slot)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *slotRepository) Update(ctx context.Context, slot domain.Slot) error {
	op := "SlotRepository.Update"
	query := `
		UPDATE slots 
		SET week_id = :week_id, tag_id = :tag_id, day_of_week = :day_of_week, activity = :activity, start_time = :start_time, end_time = :end_time, updated_at = :updated_at
		WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, slot)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *slotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "SlotRepository.Delete"
	query := `DELETE FROM slots WHERE id = $1`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
