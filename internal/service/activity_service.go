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
	ErrActivityNotFound = errors.New("activity not found")
)

type ActivityService interface {
	Get(ctx context.Context, id string) (*domain.Activity, error)
	List(ctx context.Context) []domain.Activity
	Create(ctx context.Context, title string) (string, error)
	Update(ctx context.Context, id string, title string) (string, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
	SetSlotCascadeDeleter(deleter SlotCascadeDeleter)
}

type SlotCascadeDeleter interface {
	DeleteByActivityID(ctx context.Context, activityID string) error
}

type ActivityServiceOption func(*activityService)

func WithSlotCascadeDeleter(deleter SlotCascadeDeleter) ActivityServiceOption {
	return func(s *activityService) {
		s.cascadeDeleter = deleter
	}
}

type activityService struct {
	mu             sync.RWMutex
	repo           map[string]domain.Activity
	cascadeDeleter SlotCascadeDeleter
}

func NewActivityService(opts ...ActivityServiceOption) ActivityService {
	s := &activityService{
		repo: make(map[string]domain.Activity, 32),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *activityService) Get(ctx context.Context, id string) (*domain.Activity, error) {
	op := "ActivityService.Get"

	s.mu.RLock()
	defer s.mu.RUnlock()

	activity, ok := s.repo[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrActivityNotFound)
	}

	// Return a copy to avoid external mutation of internal state
	result := activity
	return &result, nil
}

func (s *activityService) List(ctx context.Context) []domain.Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activities := make([]domain.Activity, 0, len(s.repo))
	for _, activity := range s.repo {
		activities = append(activities, activity)
	}

	return activities
}

func (s *activityService) Create(ctx context.Context, title string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	now := time.Now().UTC()

	s.repo[id] = domain.Activity{
		ID:        id,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return id, nil
}

func (s *activityService) Update(ctx context.Context, id string, title string) (string, error) {
	op := "ActivityService.Update"

	s.mu.Lock()
	defer s.mu.Unlock()

	activity, ok := s.repo[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrActivityNotFound)
	}

	activity.Title = title
	activity.UpdatedAt = time.Now().UTC()
	s.repo[id] = activity

	return id, nil
}

func (s *activityService) Delete(ctx context.Context, id string) error {
	op := "ActivityService.Delete"

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.repo[id]; !ok {
		return fmt.Errorf("%s: %w", op, ErrActivityNotFound)
	}

	delete(s.repo, id)

	if s.cascadeDeleter != nil {
		_ = s.cascadeDeleter.DeleteByActivityID(ctx, id)
	}

	return nil
}

func (s *activityService) DeleteAll(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id := range s.repo {
		if s.cascadeDeleter != nil {
			_ = s.cascadeDeleter.DeleteByActivityID(ctx, id)
		}
		delete(s.repo, id)
	}

	return nil
}

func (s *activityService) SetSlotCascadeDeleter(deleter SlotCascadeDeleter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cascadeDeleter = deleter
}
