package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmi3midd/shkvcache"
	"github.com/dmi3midd/simple-schedule/internal/domain"
	"github.com/dmi3midd/simple-schedule/internal/repository"
	"github.com/google/uuid"
)

type WeekService interface {
	// CreateWeek creates an empty week with the provided title
	CreateWeek(ctx context.Context, title string) (uuid.UUID, error)
	// GetWeek retrieves a week by its ID
	GetWeek(ctx context.Context, id uuid.UUID) (*domain.Week, error)
	// GetAllWeeks retrieves all weeks
	GetAllWeeks(ctx context.Context) ([]domain.Week, error)
	// GetWeekSchedule retrieves the full schedule for a week, with slots and tags
	GetWeekSchedule(ctx context.Context, weekId uuid.UUID) (*domain.WeekSchedule, error)
	// UpdateWeek updates an existing week's title
	UpdateWeek(ctx context.Context, id uuid.UUID, title string) (uuid.UUID, error)
	// DeleteWeek deletes a week by its ID
	DeleteWeek(ctx context.Context, id uuid.UUID) error
}

type weekService struct {
	weekRepo      repository.WeekRepository
	slotRepo      repository.SlotRepository
	scheduleCahce shkvcache.Cache[domain.WeekSchedule]
	maxWeeks      int
}

func NewWeekService(
	weekRepo repository.WeekRepository,
	slotRepo repository.SlotRepository,
	scheduleCahce shkvcache.Cache[domain.WeekSchedule],
	maxWeeks int,
) WeekService {
	return &weekService{
		weekRepo:      weekRepo,
		slotRepo:      slotRepo,
		scheduleCahce: scheduleCahce,
		maxWeeks:      maxWeeks,
	}
}

func (s *weekService) CreateWeek(ctx context.Context, title string) (uuid.UUID, error) {
	op := "WeekService.CreateWeek"

	if s.maxWeeks > 0 {
		weeks, err := s.weekRepo.GetAll(ctx)
		if err != nil {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
		if len(weeks) >= s.maxWeeks {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrMaxWeeksReached)
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()
	week := domain.Week{
		ID:        id,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.weekRepo.Create(ctx, week); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *weekService) GetWeek(ctx context.Context, id uuid.UUID) (*domain.Week, error) {
	op := "WeekService.GetWeek"
	week, err := s.weekRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoWeek) {
			return nil, fmt.Errorf("%s: %w", op, ErrWeekNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return week, nil
}

func (s *weekService) GetAllWeeks(ctx context.Context) ([]domain.Week, error) {
	op := "WeekService.GetAllWeeks"
	weeks, err := s.weekRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return weeks, nil
}

func (s *weekService) GetWeekSchedule(ctx context.Context, weekId uuid.UUID) (*domain.WeekSchedule, error) {
	op := "WeekService.GetWeekSchedule"
	if schedule, ok := s.scheduleCahce.Get(weekId.String()); ok {
		return &schedule, nil
	}

	week, err := s.weekRepo.GetById(ctx, weekId)
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

	days := make(map[domain.DayOfWeek][]domain.SlotWithTag, 7)
	for _, day := range domain.AllDaysOfWeek() {
		days[day] = make([]domain.SlotWithTag, 0)
	}

	for _, slot := range slots {
		days[slot.DayOfWeek] = append(days[slot.DayOfWeek], slot)
	}

	schedule := domain.WeekSchedule{
		Week:  *week,
		Slots: slots,
		Days:  days,
	}

	s.scheduleCahce.Set(weekId.String(), schedule, 120)

	return &schedule, nil
}

func (s *weekService) UpdateWeek(ctx context.Context, id uuid.UUID, title string) (uuid.UUID, error) {
	op := "WeekService.UpdateWeek"

	existing, err := s.weekRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoWeek) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrWeekNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	week := domain.Week{
		ID:        id,
		Title:     title,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}
	if err := s.weekRepo.Update(ctx, week); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *weekService) DeleteWeek(ctx context.Context, id uuid.UUID) error {
	op := "WeekService.DeleteWeek"
	_, err := s.weekRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoWeek) {
			return fmt.Errorf("%s: %w", op, ErrWeekNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.weekRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
