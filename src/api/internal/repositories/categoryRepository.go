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

type ICategoryRepository interface {
	WithTX(tx *gorm.DB) ICategoryRepository
	GetCategories(
		ctx context.Context,
		page, limit int,
	) ([]*models.Category, *pagination.Pagination, error)
	GetCategoryById(ctx context.Context, id interface{}, isSecondary bool) (*models.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) error
	UpdateCategory(
		ctx context.Context,
		id interface{},
		category map[string]interface{},
		isSecondary bool,
	) error
	DeleteCategory(ctx context.Context, categoryId uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (cr *categoryRepository) WithTX(tx *gorm.DB) ICategoryRepository {
	if tx == nil {
		return cr
	}
	return &categoryRepository{
		db: tx,
	}
}

func (cr *categoryRepository) GetCategories(
	ctx context.Context,
	page, limit int,
) ([]*models.Category, *pagination.Pagination, error) {
	var totalCategories int64
	if err := cr.db.WithContext(ctx).Model(&models.Category{}).Count(&totalCategories).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of categories: %w", err)
	}

	pagination := pagination.NewPagination(page, limit, int(totalCategories))

	var categories []*models.Category
	if err := cr.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return categories, pagination, nil
}

func (cr *categoryRepository) GetCategoryById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Category, error) {
	var category models.Category
	query := cr.db.WithContext(ctx)

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.Take(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		idStr := fmt.Sprintf("%v", id)
		return nil, fmt.Errorf("failed to get category with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, err)
	}
	return &category, nil
}

func (cr *categoryRepository) GetCategoryBySlug(
	ctx context.Context,
	slug string,
) (*models.Category, error) {
	var category models.Category
	if err := cr.db.WithContext(ctx).Take(&category, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (cr *categoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	if err := cr.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (cr *categoryRepository) UpdateCategory(
	ctx context.Context,
	id interface{},
	category map[string]interface{},
	isSecondary bool,
) error {
	query := cr.db.WithContext(ctx).Model(&models.Category{})

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	result := query.Updates(category)
	if result.Error != nil {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("failed to update category with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, result.Error)
	}

	if result.RowsAffected == 0 {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("category with %s %s not found",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr)
	}
	return nil
}

func (cr *categoryRepository) DeleteCategory(ctx context.Context, categoryId uuid.UUID) error {
	result := cr.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", categoryId).
		Delete(&models.Category{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete category with id %s: %w", categoryId, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("category with id %s not found", categoryId)
	}
	return nil
}
