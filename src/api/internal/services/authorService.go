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
	GetAuthors(
		ctx context.Context,
		page, limit int,
	) ([]*models.Author, *pagination.Pagination, error)
	GetAuthorById(ctx context.Context, id interface{}, isSecondary bool) (*models.Author, error)
	GetAuthorBySlug(ctx context.Context, slug string) (*models.Author, error)
	CreateAuthor(
		ctx context.Context,
		authorRequest dtos.CreateAuthorRequest,
		currentUserId uuid.UUID,
	) (*models.Author, error)
	UpdateAuthor(
		ctx context.Context,
		id interface{},
		authorRequest dtos.UpdateAuthorRequest,
		currentUserId uuid.UUID,
		isSecondary bool,
	) (*models.Author, error)
	DeleteAuthor(ctx context.Context, authorId uuid.UUID, currentUserId uuid.UUID) error
}

type authorService struct {
	db               *gorm.DB
	authorRepository repositories.IAuthorRepository
	userRepository   repositories.IUserRepository
}

func NewAuthorService(
	db *gorm.DB,
	authorRepository repositories.IAuthorRepository,
	userRepository repositories.IUserRepository,
) IAuthorService {
	return &authorService{
		db:               db,
		authorRepository: authorRepository,
		userRepository:   userRepository,
	}
}

func (as *authorService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.IAuthorRepository, repositories.IUserRepository) error,
) error {
	return as.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		authorRepository := as.authorRepository.WithTX(tx)
		userRepository := as.userRepository.WithTX(tx)
		return fn(ctx, authorRepository, userRepository)
	})
}

func (as *authorService) GetAuthors(
	ctx context.Context,
	page, limit int,
) ([]*models.Author, *pagination.Pagination, error) {
	return as.authorRepository.GetAuthors(ctx, page, limit)
}

func (as *authorService) GetAuthorById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Author, error) {
	return as.authorRepository.GetAuthorById(ctx, id, isSecondary)
}

func (as *authorService) GetAuthorBySlug(
	ctx context.Context,
	authorSlug string,
) (*models.Author, error) {
	return as.authorRepository.GetAuthorBySlug(ctx, authorSlug)
}

func (as *authorService) CreateAuthor(
	ctx context.Context,
	authorRequest dtos.CreateAuthorRequest,
	currentUserId uuid.UUID,
) (*models.Author, error) {
	var createdAuthor *models.Author
	err := as.withTX(
		ctx,
		func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
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
				SecondaryID: authorRequest.SecondaryID,
				Name:        authorRequest.Name,
				Slug:        authorSlug,
			}
			if err := authorRepository.CreateAuthor(ctx, author); err != nil {
				return err
			}

			// Get the created author with all fields populated
			createdAuthor, err = authorRepository.GetAuthorById(ctx, author.ID, false)
			return err
		},
	)
	return createdAuthor, err
}

func (as *authorService) UpdateAuthor(
	ctx context.Context,
	id interface{},
	authorRequest dtos.UpdateAuthorRequest,
	currentUserId uuid.UUID,
	isSecondary bool,
) (*models.Author, error) {
	var updatedAuthor *models.Author
	err := as.withTX(
		ctx,
		func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
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

			author, err := authorRepository.GetAuthorById(ctx, id, isSecondary)
			if err != nil {
				return err
			}
			if author == nil {
				idStr := fmt.Sprintf("%v", id)
				return fmt.Errorf("author with %s %s not found",
					map[bool]string{true: "secondary id", false: "id"}[isSecondary],
					idStr)
			}

			var authorToUpdate models.Author
			if authorRequest.SecondaryID != nil {
				authorToUpdate.SecondaryID = *authorRequest.SecondaryID
			}
			if authorRequest.Name != nil && *authorRequest.Name != author.Name {
				authorToUpdate.Name = *authorRequest.Name
				slugStr := slug.Make(*authorRequest.Name)
				if slugStr != author.Slug {
					existingAuthor, err := authorRepository.GetAuthorBySlug(ctx, slugStr)
					if err != nil {
						if existingAuthor != nil {
							return fmt.Errorf("author with slug %s already exists", slugStr)
						}
						return err
					}
					authorToUpdate.Slug = slugStr
				}
			}
			if err := authorRepository.UpdateAuthor(ctx, id, &authorToUpdate, isSecondary); err != nil {
				return err
			}

			// Get the updated author
			updatedAuthor, err = authorRepository.GetAuthorById(ctx, id, isSecondary)
			return err
		},
	)
	return updatedAuthor, err
}

func (as *authorService) DeleteAuthor(
	ctx context.Context,
	authorId uuid.UUID,
	currentUserId uuid.UUID,
) error {
	return as.withTX(
		ctx,
		func(ctx context.Context, authorRepository repositories.IAuthorRepository, userRepository repositories.IUserRepository) error {
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
		},
	)
}
