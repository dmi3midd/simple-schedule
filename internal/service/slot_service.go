package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/dmi3midd/simple-schedule/internal/domain"
)

var (
	ErrSlotNotFound = errors.New("slot not found")
)

type CreateSlotInput struct {
	ActivityID string
	Day        string
	StartTime  string
	Duration   string
}

type UpdateSlotInput struct {
	ActivityID string
	Day        string
	StartTime  string
	Duration   string
}

type SlotService interface {
	Get(ctx context.Context, id string) (*domain.Slot, error)
	List(ctx context.Context) []domain.Slot
	ListByActivityID(ctx context.Context, activityID string) []domain.Slot
	Create(ctx context.Context, input CreateSlotInput) (string, error)
	Update(ctx context.Context, id string, input UpdateSlotInput) (string, error)
	Delete(ctx context.Context, id string) error
	DeleteByActivityID(ctx context.Context, activityID string) error
}

type slotService struct {
	mu              sync.RWMutex
	repo            map[string]domain.Slot
	activityService ActivityService
}

func NewSlotService(activityService ActivityService) SlotService {
	return &slotService{
		repo:            make(map[string]domain.Slot, 32),
		activityService: activityService,
	}
}

func (s *slotService) Get(ctx context.Context, id string) (*domain.Slot, error) {
	op := "SlotService.Get"

	s.mu.RLock()
	defer s.mu.RUnlock()

	slot, ok := s.repo[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrSlotNotFound)
	}

	result := slot
	return &result, nil
}

func (s *slotService) List(ctx context.Context) []domain.Slot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slots := make([]domain.Slot, 0, len(s.repo))
	for _, slot := range s.repo {
		slots = append(slots, slot)
	}

	return slots
}

func (s *slotService) ListByActivityID(ctx context.Context, activityID string) []domain.Slot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slots := make([]domain.Slot, 0)
	for _, slot := range s.repo {
		if slot.ActivityID == activityID {
			slots = append(slots, slot)
		}
	}

	return slots
}

func (s *slotService) Create(ctx context.Context, input CreateSlotInput) (string, error) {
	op := "SlotService.Create"

	if s.activityService != nil && input.ActivityID != "" {
		if _, err := s.activityService.Get(ctx, input.ActivityID); err != nil {
			return "", fmt.Errorf("%s: activity check failed: %w", op, err)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	now := time.Now().UTC()

	s.repo[id] = domain.Slot{
		ID:         id,
		ActivityID: input.ActivityID,
		Day:        input.Day,
		StartTime:  input.StartTime,
		Duration:   input.Duration,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return id, nil
}

func (s *slotService) Update(ctx context.Context, id string, input UpdateSlotInput) (string, error) {
	op := "SlotService.Update"

	if input.ActivityID != "" && s.activityService != nil {
		if _, err := s.activityService.Get(ctx, input.ActivityID); err != nil {
			return "", fmt.Errorf("%s: activity check failed: %w", op, err)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	slot, ok := s.repo[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrSlotNotFound)
	}

	slot.ActivityID = input.ActivityID
	slot.Day = input.Day
	slot.StartTime = input.StartTime
	slot.Duration = input.Duration
	slot.UpdatedAt = time.Now().UTC()

	s.repo[id] = slot

	return id, nil
}

func (s *slotService) Delete(ctx context.Context, id string) error {
	op := "SlotService.Delete"

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.repo[id]; !ok {
		return fmt.Errorf("%s: %w", op, ErrSlotNotFound)
	}

	delete(s.repo, id)
	return nil
}

func (s *slotService) DeleteByActivityID(ctx context.Context, activityID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, slot := range s.repo {
		if slot.ActivityID == activityID {
			delete(s.repo, id)
		}
	}

	return nil
}
