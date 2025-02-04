package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type ICategoryService interface {
	GetCategories(ctx context.Context, page, limit int) ([]*models.Category, *pagination.Pagination, error)
	GetProductById(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
	CreateCategory(ctx context.Context, categoryRequest dtos.CreateCategoryRequest, currentUserId uuid.UUID) error
	UpdateCategory(ctx context.Context, categoryId uuid.UUID, categoryRequest dtos.UpdateCategoryRequest, currentUserId uuid.UUID) error
	DeleteCategory(ctx context.Context, categoryId uuid.UUID, currentUserId uuid.UUID) error
}

type categoryService struct {
	db                 *gorm.DB
	categoryRepository repositories.ICategoryRepository
	userRepository     repositories.IUserRepository
}

func NewCategoryService(db *gorm.DB, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) ICategoryService {
	return &categoryService{
		db:                 db,
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
	}
}

func (cs *categoryService) withTx(ctx context.Context, fn func(context.Context, repositories.ICategoryRepository, repositories.IUserRepository) error) error {
	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		categoryRepository := cs.categoryRepository.WithTX(tx)
		userRepository := cs.userRepository.WithTX(tx)
		return fn(ctx, categoryRepository, userRepository)
	})
}

func (cs *categoryService) GetCategories(ctx context.Context, page, limit int) ([]*models.Category, *pagination.Pagination, error) {
	return cs.categoryRepository.GetCategories(ctx, page, limit)
}

func (cs *categoryService) GetProductById(ctx context.Context, categoryId uuid.UUID) (*models.Category, error) {
	return cs.categoryRepository.GetProductById(ctx, categoryId)
}

func (cs *categoryService) CreateCategory(ctx context.Context, categoryRequest dtos.CreateCategoryRequest, currentUserId uuid.UUID) error {
	return cs.withTx(ctx, func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
		// Check if the user is an admin
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can create category")
		}

		// Construct a slug for the category
		slug := slug.Make(categoryRequest.Name)

		// Check if the category already exists
		existingCategory, err := categoryRepository.GetCategoryBySlug(ctx, slug)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingCategory != nil {
			return errors.New("category already exists")
		}

		// Define a new category
		category := &models.Category{
			Name: categoryRequest.Name,
			Slug: slug,
		}

		return categoryRepository.CreateCategory(ctx, category)
	})
}

func (cs *categoryService) UpdateCategory(ctx context.Context, categoryId uuid.UUID, categoryRequest dtos.UpdateCategoryRequest, currentUserId uuid.UUID) error {
	return cs.withTx(ctx, func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
		// Check if the user is an admin
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can update category")
		}

		// Get the category by id
		category, err := categoryRepository.GetProductById(ctx, categoryId)
		if err != nil {
			return err
		}

		// Apply the updates
		if categoryRequest.Name != nil && *categoryRequest.Name != category.Name {
			category.Name = *categoryRequest.Name
			slug := slug.Make(*categoryRequest.Name)

			// Check if the category already exists
			existingCategory, err := categoryRepository.GetCategoryBySlug(ctx, slug)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if existingCategory != nil {
				return errors.New("category already exists")
			}

			category.Slug = slug
		}

		return categoryRepository.UpdateCategory(ctx, categoryId, category)
	})
}

func (cs *categoryService) DeleteCategory(ctx context.Context, categoryId uuid.UUID, currentUserId uuid.UUID) error {
	return cs.withTx(ctx, func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
		// Check if the user is an admin
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can delete category")
		}

		return categoryRepository.DeleteCategory(ctx, categoryId)
	})
}
