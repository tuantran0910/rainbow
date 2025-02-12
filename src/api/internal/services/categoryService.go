package services

import (
	"context"
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
	GetCategories(
		ctx context.Context,
		page, limit int,
	) ([]*models.Category, *pagination.Pagination, error)
	GetCategoryById(ctx context.Context, categoryId uuid.UUID) (*models.Category, error)
	GetCategoryBySlug(ctx context.Context, categorySlug string) (*models.Category, error)
	CreateCategory(
		ctx context.Context,
		categoryRequest dtos.CreateCategoryRequest,
		currentUserId uuid.UUID,
	) error
	UpdateCategory(
		ctx context.Context,
		categoryId uuid.UUID,
		categoryRequest dtos.UpdateCategoryRequest,
		currentUserId uuid.UUID,
	) error
	DeleteCategory(ctx context.Context, categoryId uuid.UUID, currentUserId uuid.UUID) error
}

type categoryService struct {
	db                 *gorm.DB
	categoryRepository repositories.ICategoryRepository
	userRepository     repositories.IUserRepository
}

func NewCategoryService(
	db *gorm.DB,
	categoryRepository repositories.ICategoryRepository,
	userRepository repositories.IUserRepository,
) ICategoryService {
	return &categoryService{
		db:                 db,
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
	}
}

func (cs *categoryService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.ICategoryRepository, repositories.IUserRepository) error,
) error {
	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		categoryRepository := cs.categoryRepository.WithTX(tx)
		userRepository := cs.userRepository.WithTX(tx)
		return fn(ctx, categoryRepository, userRepository)
	})
}

func (cs *categoryService) GetCategories(
	ctx context.Context,
	page, limit int,
) ([]*models.Category, *pagination.Pagination, error) {
	return cs.categoryRepository.GetCategories(ctx, page, limit)
}

func (cs *categoryService) GetCategoryById(
	ctx context.Context,
	categoryId uuid.UUID,
) (*models.Category, error) {
	return cs.categoryRepository.GetCategoryById(ctx, categoryId)
}

func (cs *categoryService) GetCategoryBySlug(
	ctx context.Context,
	categorySlug string,
) (*models.Category, error) {
	return cs.categoryRepository.GetCategoryBySlug(ctx, categorySlug)
}

func (cs *categoryService) CreateCategory(
	ctx context.Context,
	categoryRequest dtos.CreateCategoryRequest,
	currentUserId uuid.UUID,
) error {
	return cs.withTX(
		ctx,
		func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			if currentUser.Role != models.AdminRole {
				return fmt.Errorf("current user does not have permission to create category")
			}

			categorySlug := categoryRequest.Slug
			if categorySlug == "" {
				categorySlug = slug.Make(categoryRequest.Name)
			}

			existingCategory, err := categoryRepository.GetCategoryBySlug(ctx, categorySlug)
			if err != nil {
				if existingCategory != nil {
					return fmt.Errorf("category with slug %s already existed", categorySlug)
				}
				return err
			}

			category := &models.Category{
				Name: categoryRequest.Name,
				Slug: categorySlug,
			}
			return categoryRepository.CreateCategory(ctx, category)
		},
	)
}

func (cs *categoryService) UpdateCategory(
	ctx context.Context,
	categoryId uuid.UUID,
	categoryRequest dtos.UpdateCategoryRequest,
	currentUserId uuid.UUID,
) error {
	return cs.withTX(
		ctx,
		func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			if currentUser.Role != models.AdminRole {
				return fmt.Errorf("current user does not have permission to update ")
			}

			category, err := categoryRepository.GetCategoryById(ctx, categoryId)
			if err != nil {
				return err
			}
			if category == nil {
				return fmt.Errorf("category with id %s not found", categoryId)
			}

			if categoryRequest.Name != nil && *categoryRequest.Name != category.Name {
				category.Name = *categoryRequest.Name
			}
			return categoryRepository.UpdateCategory(ctx, categoryId, category)
		},
	)
}

func (cs *categoryService) DeleteCategory(
	ctx context.Context,
	categoryId uuid.UUID,
	currentUserId uuid.UUID,
) error {
	return cs.withTX(
		ctx,
		func(ctx context.Context, categoryRepository repositories.ICategoryRepository, userRepository repositories.IUserRepository) error {
			// Check if the user is an admin
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			if currentUser.Role != models.AdminRole {
				return fmt.Errorf("current user does not have permission to delete category")
			}
			return categoryRepository.DeleteCategory(ctx, categoryId)
		},
	)
}
