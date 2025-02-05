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

type IAuthorService interface {
	GetAuthors(ctx context.Context, page, limit int) ([]*models.Author, *pagination.Pagination, error)
	GetAuthorById(ctx context.Context, authorId uuid.UUID) (*models.Author, error)
	GetAuthorBySlug(ctx context.Context, authorSlug string) (*models.Author, error)
	CreateAuthor(ctx context.Context, authorRequest dtos.CreateAuthorRequest, currentUserId uuid.UUID) error
	UpdateAuthor(ctx context.Context, authorId uuid.UUID, authorRequest dtos.UpdateAuthorRequest, currentUserId uuid.UUID) error
	DeleteAuthor(ctx context.Context, authorId uuid.UUID, currentUserId uuid.UUID) error
}

type authorService struct {
	db               *gorm.DB
	authorRepository repositories.IAuthorRepository
	userRepository   repositories.IUserRepository
}

func NewAuthorService(db *gorm.DB, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) IAuthorService {
	return &authorService{
		db:               db,
		authorRepository: authorRepository,
		userRepository:   userRepository,
	}
}

func (as *authorService) withTx(ctx context.Context, fn func(context.Context, repositories.IAuthorRepository, repositories.IUserRepository) error) error {
	return as.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		authorRepository := as.authorRepository.WithTX(tx)
		userRepository := as.userRepository.WithTX(tx)
		return fn(ctx, authorRepository, userRepository)
	})
}

func (as *authorService) GetAuthors(ctx context.Context, page, limit int) ([]*models.Author, *pagination.Pagination, error) {
	return as.authorRepository.GetAuthors(ctx, page, limit)
}

func (as *authorService) GetAuthorById(ctx context.Context, authorId uuid.UUID) (*models.Author, error) {
	return as.authorRepository.GetAuthorById(ctx, authorId)
}

func (as *authorService) GetAuthorBySlug(ctx context.Context, authorSlug string) (*models.Author, error) {
	return as.authorRepository.GetAuthorBySlug(ctx, authorSlug)
}

func (as *authorService) CreateAuthor(ctx context.Context, authorRequest dtos.CreateAuthorRequest, currentUserId uuid.UUID) error {
	return as.withTx(ctx, func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.Role != models.AdminRole {
			return fmt.Errorf("current user does not have permission to create author")
		}

		authorSlug := authorRequest.Slug
		if authorSlug == "" {
			authorSlug = slug.Make(authorRequest.Name)
		}

		existingAuthor, err := authorRepository.GetAuthorBySlug(ctx, authorSlug)
		if err != nil {
			if existingAuthor != nil {
				return fmt.Errorf("author with slug %s already existed", authorSlug)
			}
			return err
		}

		author := &models.Author{
			Name: authorRequest.Name,
			Slug: authorSlug,
		}
		return authorRepository.CreateAuthor(ctx, author)
	})
}

func (as *authorService) UpdateAuthor(ctx context.Context, authorId uuid.UUID, authorRequest dtos.UpdateAuthorRequest, currentUserId uuid.UUID) error {
	return as.withTx(ctx, func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.Role != models.AdminRole {
			return fmt.Errorf("current user does not have permission to update author")
		}

		author, err := authorRepository.GetAuthorById(ctx, authorId)
		if err != nil {
			return err
		}
		if author == nil {
			return fmt.Errorf("author with id %s not found", authorId)
		}

		if authorRequest.Name != nil && *authorRequest.Name != author.Name {
			author.Name = *authorRequest.Name
		}
		return authorRepository.UpdateAuthor(ctx, authorId, author)
	})
}

func (as *authorService) DeleteAuthor(ctx context.Context, authorId uuid.UUID, currentUserId uuid.UUID) error {
	return as.withTx(ctx, func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.Role != models.AdminRole {
			return fmt.Errorf("current user does not have permission to delete author")
		}
		return authorRepository.DeleteAuthor(ctx, authorId)
	})
}
