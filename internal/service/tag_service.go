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

type TagService interface {
	// GetAll returns all tags
	GetAll(ctx context.Context) ([]domain.Tag, error)
	// GetById returns a tag by its ID
	GetById(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	// Create creates a new tag
	Create(ctx context.Context, title, hexColor string) (uuid.UUID, error)
	// Update updates an existing tag
	Update(ctx context.Context, id uuid.UUID, title, hexColor string) (uuid.UUID, error)
	// Delete deletes a tag
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) GetAll(ctx context.Context) ([]domain.Tag, error) {
	op := "TagService.GetAll"
	tags, err := s.tagRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tags, nil
}

func (s *tagService) GetById(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	op := "TagService.GetById"
	tag, err := s.tagRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTag) {
			return nil, fmt.Errorf("%s: %w", op, ErrTagNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tag, nil
}

func (s *tagService) Create(ctx context.Context, title, hexColor string) (uuid.UUID, error) {
	op := "TagService.Create"
	existing, err := s.tagRepo.GetByTitle(ctx, title)
	if err != nil && !errors.Is(err, repository.ErrNoTag) {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	if existing != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, ErrTagAlreadyExists)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	now := time.Now()
	tag := domain.Tag{
		ID:        id,
		Title:     title,
		HexColor:  hexColor,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *tagService) Update(ctx context.Context, id uuid.UUID, title, hexColor string) (uuid.UUID, error) {
	op := "TagService.Update"
	existing, err := s.tagRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTag) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrTagNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	if title != existing.Title {
		conflict, err := s.tagRepo.GetByTitle(ctx, title)
		if err != nil && !errors.Is(err, repository.ErrNoTag) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, err)
		}
		if conflict != nil {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrTagAlreadyExists)
		}
	}

	now := time.Now()
	tag := domain.Tag{
		ID:        id,
		Title:     title,
		HexColor:  hexColor,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: now,
	}
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	op := "TagService.Delete"
	_, err := s.tagRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTag) {
			return fmt.Errorf("%s: %w", op, ErrTagNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.tagRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
