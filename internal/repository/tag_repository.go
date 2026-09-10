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
	ErrNoTag = errors.New("no tag found")
)

type TagRepository interface {
	// GetById retrieves a tag by its ID.
	// Returns [ErrNoTag] if the tag is not found.
	GetById(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	// GetAll retrieves all tags
	// Returns an empty slice if no tags are found.
	GetAll(ctx context.Context) ([]domain.Tag, error)
	// Create creates a new tag
	Create(ctx context.Context, tag domain.Tag) (*domain.Tag, error)
	// Update updates an existing tag
	Update(ctx context.Context, tag domain.Tag) (*domain.Tag, error)
	// Delete deletes a tag
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	op := "TagRepository.GetById"
	query := "SELECT id, title, hex_color, created_at, updated_at FROM tags WHERE id = $1"
	var tag domain.Tag
	executor := ExtractTx(ctx, r.db)
	err := sqlx.GetContext(ctx, executor, &tag, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoTag)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &tag, nil
}

func (r *tagRepository) GetAll(ctx context.Context) ([]domain.Tag, error) {
	op := "TagRepository.GetAll"
	query := `
	SELECT id, title, hex_color, created_at, updated_at 
	FROM tags
	`
	var tags []domain.Tag
	executor := ExtractTx(ctx, r.db)
	err := sqlx.SelectContext(ctx, executor, &tags, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tags, nil
}

func (r *tagRepository) Create(ctx context.Context, tag domain.Tag) (*domain.Tag, error) {
	op := "TagRepository.Create"
	query := `
	INSERT INTO tags (id, title, hex_color, created_at, updated_at) 
	VALUES (:id, :title, :hex_color, :created_at, :updated_at) 
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, tag)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &tag, nil
}

func (r *tagRepository) Update(ctx context.Context, tag domain.Tag) (*domain.Tag, error) {
	op := "TagRepository.Update"
	query := `
	UPDATE tags 
	SET title = :title, hex_color = :hex_color, updated_at = :updated_at
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, tag)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &tag, nil
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "TagRepository.Delete"
	query := `
	DELETE FROM tags 
	WHERE id = :id
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
