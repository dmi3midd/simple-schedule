package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/repository"
	"github.com/google/uuid"
)

type CreateSlotInput struct {
	WeekID    uuid.UUID
	TagID     *uuid.UUID
	DayOfWeek domain.DayOfWeek
	Activity  string
	StartTime int
	EndTime   int
}

type UpdateSlotInput struct {
	ID        uuid.UUID
	WeekID    uuid.UUID
	TagID     *uuid.UUID
	DayOfWeek domain.DayOfWeek
	Activity  string
	StartTime int
	EndTime   int
}

type SlotService interface {
	// CreateSlot creates a new slot after validating time range, day, and overlaps
	CreateSlot(ctx context.Context, input CreateSlotInput) (uuid.UUID, error)
	// UpdateSlot updates an existing slot with validation
	UpdateSlot(ctx context.Context, input UpdateSlotInput) error
	// DeleteSlot deletes a slot by ID
	DeleteSlot(ctx context.Context, id uuid.UUID) error
	// GetSlot retrieves a slot by ID with its associated tag
	GetSlot(ctx context.Context, id uuid.UUID) (*domain.SlotWithTag, error)
	// GetSlotsByWeek retrieves all slots with tags for a week
	GetSlotsByWeek(ctx context.Context, weekId uuid.UUID) ([]domain.SlotWithTag, error)
}

type slotService struct {
	slotRepo       repository.SlotRepository
	weekRepo       repository.WeekRepository
	tagRepo        repository.TagRepository
	maxSlotsPerDay int
}

func NewSlotService(
	slotRepo repository.SlotRepository,
	weekRepo repository.WeekRepository,
	tagRepo repository.TagRepository,
	maxSlotsPerDay int,
) SlotService {
	return &slotService{
		slotRepo:       slotRepo,
		weekRepo:       weekRepo,
		tagRepo:        tagRepo,
		maxSlotsPerDay: maxSlotsPerDay,
	}
}

func (s *slotService) validateSlot(ctx context.Context, weekID uuid.UUID, tagID *uuid.UUID, dayOfWeek domain.DayOfWeek, startTime, endTime int, excludeSlotID *uuid.UUID) error {
	if startTime < 0 || endTime > 86400 || startTime >= endTime {
		return ErrInvalidTimeRange
	}

	if !dayOfWeek.IsValid() {
		return ErrInvalidDayOfWeek
	}

	_, err := s.weekRepo.GetById(ctx, weekID)
	if err != nil {
		if errors.Is(err, repository.ErrNoWeek) {
			return ErrWeekNotFound
		}
		return err
	}

	if tagID != nil {
		_, err := s.tagRepo.GetById(ctx, *tagID)
		if err != nil {
			if errors.Is(err, repository.ErrNoTag) {
				return ErrTagNotFound
			}
			return err
		}
	}

	if s.maxSlotsPerDay > 0 {
		slots, err := s.slotRepo.GetByWeekId(ctx, weekID)
		if err != nil {
			return err
		}
		daySlotsCount := 0
		for _, slot := range slots {
			if excludeSlotID != nil && slot.ID == *excludeSlotID {
				continue
			}
			if slot.DayOfWeek == dayOfWeek {
				daySlotsCount++
			}
		}
		if daySlotsCount >= s.maxSlotsPerDay {
			return ErrMaxSlotsPerDayReached
		}
	}

	overlap, err := s.slotRepo.HasOverlap(ctx, weekID, dayOfWeek, startTime, endTime, excludeSlotID)
	if err != nil {
		return err
	}
	if overlap {
		return ErrSlotOverlap
	}

	return nil
}

func (s *slotService) CreateSlot(ctx context.Context, input CreateSlotInput) (uuid.UUID, error) {
	op := "SlotService.CreateSlot"

	if err := s.validateSlot(ctx, input.WeekID, input.TagID, input.DayOfWeek, input.StartTime, input.EndTime, nil); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	slot := domain.Slot{
		ID:        id,
		WeekID:    input.WeekID,
		TagID:     input.TagID,
		DayOfWeek: input.DayOfWeek,
		Activity:  input.Activity,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.slotRepo.Create(ctx, slot); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *slotService) UpdateSlot(ctx context.Context, input UpdateSlotInput) error {
	op := "SlotService.UpdateSlot"

	existing, err := s.slotRepo.GetById(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNoSlot) {
			return fmt.Errorf("%s: %w", op, ErrSlotNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.validateSlot(ctx, input.WeekID, input.TagID, input.DayOfWeek, input.StartTime, input.EndTime, &input.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	slot := domain.Slot{
		ID:        input.ID,
		WeekID:    input.WeekID,
		TagID:     input.TagID,
		DayOfWeek: input.DayOfWeek,
		Activity:  input.Activity,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.slotRepo.Update(ctx, slot); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *slotService) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	op := "SlotService.DeleteSlot"

	_, err := s.slotRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSlot) {
			return fmt.Errorf("%s: %w", op, ErrSlotNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.slotRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *slotService) GetSlot(ctx context.Context, id uuid.UUID) (*domain.SlotWithTag, error) {
	op := "SlotService.GetSlot"
	slot, err := s.slotRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSlot) {
			return nil, fmt.Errorf("%s: %w", op, ErrSlotNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return slot, nil
}

func (s *slotService) GetSlotsByWeek(ctx context.Context, weekId uuid.UUID) ([]domain.SlotWithTag, error) {
	op := "SlotService.GetSlotsByWeek"

	_, err := s.weekRepo.GetById(ctx, weekId)
	if err != nil {
		if errors.Is(err, repository.ErrNoWeek) {
			return nil, fmt.Errorf("%s: %w", op, ErrWeekNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slots, err := s.slotRepo.GetByWeekId(ctx, weekId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return slots, nil
}
