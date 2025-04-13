package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IBookService interface {
	GetBooks(ctx context.Context, page, limit int) ([]*models.Book, *pagination.Pagination, error)
	GetBookById(ctx context.Context, id interface{}, isSecondary bool) (*models.Book, error)
	CreateBook(
		ctx context.Context,
		bookRequest dtos.CreateBookRequest,
		currentUserId uuid.UUID,
	) (*models.Book, error)
	UpdateBook(
		ctx context.Context,
		bookId uuid.UUID,
		bookRequest dtos.UpdateBookRequest,
		currentUserId uuid.UUID,
	) (*models.Book, error)
	DeleteBook(ctx context.Context, bookId uuid.UUID, currentUserId uuid.UUID) error
	CountBooks(ctx context.Context) (int64, error)
}

type bookService struct {
	db                  *gorm.DB
	bookRepository      repositories.IBookRepository
	inventoryRepository repositories.IInventoryRepository
	sellerRepository    repositories.ISellerRepository
	userRepository      repositories.IUserRepository
}

func NewBookService(
	db *gorm.DB,
	bookRepository repositories.IBookRepository,
	inventoryRepository repositories.IInventoryRepository,
	sellerRepository repositories.ISellerRepository,
	userRepository repositories.IUserRepository,
) IBookService {
	return &bookService{
		db:                  db,
		bookRepository:      bookRepository,
		inventoryRepository: inventoryRepository,
		sellerRepository:    sellerRepository,
		userRepository:      userRepository,
	}
}

func (ps *bookService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.IBookRepository, repositories.IInventoryRepository, repositories.ISellerRepository, repositories.IUserRepository) error,
) error {
	return ps.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bookRepository := ps.bookRepository.WithTX(tx)
		inventoryRepository := ps.inventoryRepository.WithTX(tx)
		sellerRepository := ps.sellerRepository.WithTX(tx)
		userRepository := ps.userRepository.WithTX(tx)
		return fn(ctx, bookRepository, inventoryRepository, sellerRepository, userRepository)
	})
}

func (ps *bookService) GetBooks(
	ctx context.Context,
	page, limit int,
) ([]*models.Book, *pagination.Pagination, error) {
	return ps.bookRepository.GetBooks(ctx, page, limit)
}

func (ps *bookService) GetBookById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Book, error) {
	return ps.bookRepository.GetBookById(ctx, id, isSecondary)
}

func (ps *bookService) CreateBook(
	ctx context.Context,
	bookRequest dtos.CreateBookRequest,
	currentUserId uuid.UUID,
) (*models.Book, error) {
	var createdBook *models.Book
	err := ps.withTX(
		ctx,
		func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			book := &models.Book{
				SecondaryID:   bookRequest.SecondaryID,
				CategoryID:    bookRequest.CategoryID,
				SellerID:      bookRequest.SellerID,
				Name:          bookRequest.Name,
				Description:   bookRequest.Description,
				Price:         bookRequest.Price,
				OriginalPrice: bookRequest.OriginalPrice,
				RatingAverage: bookRequest.RatingAverage,
				ReviewCount:   bookRequest.ReviewCount,
				PageCount:     bookRequest.PageCount,
				SoldCount:     0,
			}
			if err := bookRepository.CreateBook(ctx, book); err != nil {
				return err
			}

			bookInventory := &models.Inventory{
				BookID:          book.ID,
				Stock:           bookRequest.Stock,
				LastRestockedAt: time.Now(),
			}
			if err := inventoryRepository.CreateInventory(ctx, bookInventory); err != nil {
				return err
			}

			for _, authorId := range bookRequest.AuthorIds {
				bookAuthor := &models.BookAuthor{
					BookID:   book.ID,
					AuthorID: authorId,
				}
				if err := bookRepository.AddAuthorToBook(ctx, bookAuthor); err != nil {
					return err
				}
			}

			// Get the created book with all fields populated
			createdBook, err = bookRepository.GetBookById(ctx, book.ID, false)
			return err
		},
	)
	return createdBook, err
}

func (ps *bookService) UpdateBook(
	ctx context.Context,
	bookId uuid.UUID,
	bookRequest dtos.UpdateBookRequest,
	currentUserId uuid.UUID,
) (*models.Book, error) {
	var updatedBook *models.Book
	err := ps.withTX(
		ctx,
		func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			// Check if user exists
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			// For now, we'll skip the seller check since we're removing the user_id from sellers
			// This will need to be updated with a new approach for seller identification

			// Get existing book
			book, err := bookRepository.GetBookById(ctx, bookId, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book not found")
			}

			// Update book fields
			if bookRequest.Name != nil {
				book.Name = *bookRequest.Name
			}
			if bookRequest.Description != nil {
				book.Description = *bookRequest.Description
			}
			if bookRequest.Price != nil {
				book.Price = *bookRequest.Price
			}
			if bookRequest.OriginalPrice != nil {
				book.OriginalPrice = *bookRequest.OriginalPrice
			}
			if bookRequest.RatingAverage != nil {
				book.RatingAverage = *bookRequest.RatingAverage
			}
			if bookRequest.ReviewCount != nil {
				book.ReviewCount = *bookRequest.ReviewCount
			}
			if bookRequest.PageCount != nil {
				book.PageCount = *bookRequest.PageCount
			}

			// Update book
			if err := bookRepository.UpdateBook(ctx, bookId, book); err != nil {
				return err
			}

			// Update inventory if stock is provided
			if bookRequest.Stock != nil {
				inventory, err := inventoryRepository.GetInventoryByBookID(ctx, bookId)
				if err != nil {
					return err
				}
				if inventory == nil {
					return fmt.Errorf("inventory not found")
				}

				inventory.Stock = *bookRequest.Stock
				inventory.LastRestockedAt = time.Now()

				if err := inventoryRepository.UpdateInventory(ctx, inventory.ID, inventory); err != nil {
					return err
				}
			}

			// Get the updated book
			updatedBook, err = bookRepository.GetBookById(ctx, bookId, false)
			return err
		},
	)

	if err != nil {
		return nil, err
	}

	return updatedBook, nil
}

func (ps *bookService) DeleteBook(
	ctx context.Context,
	bookId uuid.UUID,
	currentUserId uuid.UUID,
) error {
	return ps.withTX(
		ctx,
		func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			// Check if user exists
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			// For now, we'll skip the seller check since we're removing the user_id from sellers
			// This will need to be updated with a new approach for seller identification

			// Get existing book
			book, err := bookRepository.GetBookById(ctx, bookId, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book not found")
			}

			// Delete book
			return bookRepository.DeleteBook(ctx, bookId)
		},
	)
}

func (ps *bookService) CountBooks(ctx context.Context) (int64, error) {
	return ps.bookRepository.CountBooks(ctx)
}
