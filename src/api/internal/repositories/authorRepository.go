package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IAuthorRepository interface {
	WithTX(tx *gorm.DB) IAuthorRepository
	GetAuthors(
		ctx context.Context,
		page, limit int,
	) ([]*models.Author, *pagination.Pagination, error)
	GetAuthorById(ctx context.Context, id interface{}, isSecondary bool) (*models.Author, error)
	GetAuthorBySlug(ctx context.Context, slug string) (*models.Author, error)
	CreateAuthor(ctx context.Context, author *models.Author) error
	UpdateAuthor(
		ctx context.Context,
		id interface{},
		author map[string]interface{},
		isSecondary bool,
	) error
	DeleteAuthor(ctx context.Context, authorId uuid.UUID) error
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) IAuthorRepository {
	return &authorRepository{
		db: db,
	}
}

func (ar *authorRepository) WithTX(tx *gorm.DB) IAuthorRepository {
	if tx == nil {
		return ar
	}
	return &authorRepository{
		db: tx,
	}
}

func (ar *authorRepository) GetAuthors(
	ctx context.Context,
	page, limit int,
) ([]*models.Author, *pagination.Pagination, error) {
	var totalAuthors int64
	if err := ar.db.WithContext(ctx).Model(&models.Author{}).Count(&totalAuthors).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of authors: %w", err)
	}

	pagination := pagination.NewPagination(page, limit, int(totalAuthors))

	var authors []*models.Author
	if err := ar.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch authors: %w", err)
	}
	return authors, pagination, nil
}

func (ar *authorRepository) GetAuthorById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Author, error) {
	var author models.Author
	query := ar.db.WithContext(ctx)

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.Take(&author).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		idStr := fmt.Sprintf("%v", id)
		return nil, fmt.Errorf("failed to get author with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, err)
	}
	return &author, nil
}

func (ar *authorRepository) GetAuthorBySlug(
	ctx context.Context,
	slug string,
) (*models.Author, error) {
	var author models.Author
	if err := ar.db.WithContext(ctx).Take(&author, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	return &author, nil
}

func (ar *authorRepository) CreateAuthor(ctx context.Context, author *models.Author) error {
	if err := ar.db.WithContext(ctx).Create(author).Error; err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	return nil
}

func (ar *authorRepository) UpdateAuthor(
	ctx context.Context,
	id interface{},
	author map[string]interface{},
	isSecondary bool,
) error {
	query := ar.db.WithContext(ctx).Model(&models.Author{})

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	result := query.Updates(author)
	if result.Error != nil {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("failed to update author with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, result.Error)
	}

	if result.RowsAffected == 0 {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("author with %s %s not found",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr)
	}
	return nil
}

func (ar *authorRepository) DeleteAuthor(ctx context.Context, authorId uuid.UUID) error {
	result := ar.db.WithContext(ctx).Unscoped().Where("id = ?", authorId).Delete(&models.Author{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete author with id %s: %w", authorId, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("author with id %s not found", authorId)
	}
	return nil
}
