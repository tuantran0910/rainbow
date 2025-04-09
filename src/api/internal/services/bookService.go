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

			seller, err := sellerRepository.GetSellerByUserId(ctx, currentUser.ID)
			if err != nil {
				return err
			}
			if seller == nil {
				return fmt.Errorf("seller not found")
			}

			book := &models.Book{
				SecondaryID:   bookRequest.SecondaryID,
				CategoryID:    bookRequest.CategoryID,
				SellerID:      seller.ID,
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
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			seller, err := sellerRepository.GetSellerByUserId(ctx, currentUser.ID)
			if err != nil {
				return err
			}
			if seller == nil {
				return fmt.Errorf("seller not found")
			}

			book, err := ps.bookRepository.GetBookById(ctx, bookId, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book with id %s not found", bookId)
			}

			if book.SellerID != seller.ID {
				return fmt.Errorf(
					"seller id %s does not match with current user seller id %s",
					book.SellerID,
					seller.ID,
				)
			}

			if bookRequest.SecondaryID != nil {
				book.SecondaryID = *bookRequest.SecondaryID
			}
			if bookRequest.CategoryID != nil && *bookRequest.CategoryID != book.CategoryID {
				book.CategoryID = *bookRequest.CategoryID
			}
			if bookRequest.Name != nil && *bookRequest.Name != book.Name {
				book.Name = *bookRequest.Name
			}
			if bookRequest.Description != nil && *bookRequest.Description != book.Description {
				book.Description = *bookRequest.Description
			}
			if bookRequest.Price != nil && *bookRequest.Price != book.Price {
				book.Price = *bookRequest.Price
			}
			if bookRequest.OriginalPrice != nil &&
				*bookRequest.OriginalPrice != book.OriginalPrice {
				book.OriginalPrice = *bookRequest.OriginalPrice
			}
			if bookRequest.RatingAverage != nil &&
				*bookRequest.RatingAverage != book.RatingAverage {
				book.RatingAverage = *bookRequest.RatingAverage
			}
			if bookRequest.ReviewCount != nil && *bookRequest.ReviewCount != book.ReviewCount {
				book.ReviewCount = *bookRequest.ReviewCount
			}
			if bookRequest.PageCount != nil && *bookRequest.PageCount != book.PageCount {
				book.PageCount = *bookRequest.PageCount
			}
			if err := bookRepository.UpdateBook(ctx, bookId, book); err != nil {
				return err
			}

			if bookRequest.Stock != nil {
				bookInventory, err := inventoryRepository.GetInventoryByBookID(ctx, book.ID)
				if err != nil {
					return err
				}
				if bookInventory == nil {
					return fmt.Errorf("inventory for book with id %s not found", bookId)
				}

				if *bookRequest.Stock > bookInventory.Stock {
					bookInventory.LastRestockedAt = time.Now()
				}

				bookInventory.Stock = *bookRequest.Stock
				if err := inventoryRepository.UpdateInventory(ctx, bookInventory.ID, bookInventory); err != nil {
					return err
				}
			}

			// Handle author updates if provided
			if len(bookRequest.AuthorIds) > 0 {
				// Clear existing book authors
				if err := bookRepository.ClearBookAuthors(ctx, bookId); err != nil {
					return err
				}

				// Add new authors
				for _, authorId := range bookRequest.AuthorIds {
					bookAuthor := &models.BookAuthor{
						BookID:   book.ID,
						AuthorID: authorId,
					}
					if err := bookRepository.AddAuthorToBook(ctx, bookAuthor); err != nil {
						return err
					}
				}
			}

			// Get the updated book with all fields populated
			updatedBook, err = bookRepository.GetBookById(ctx, bookId, false)
			return err
		},
	)
	return updatedBook, err
}

func (ps *bookService) DeleteBook(
	ctx context.Context,
	bookId uuid.UUID,
	currentUserId uuid.UUID,
) error {
	return ps.withTX(
		ctx,
		func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			currentUser, err := userRepository.GetUserById(ctx, currentUserId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			seller, err := sellerRepository.GetSellerByUserId(ctx, currentUser.ID)
			if err != nil {
				return err
			}
			if seller == nil {
				return fmt.Errorf("seller not found")
			}

			book, err := ps.bookRepository.GetBookById(ctx, bookId, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book with id %s not found", bookId)
			}

			if currentUser.Role != models.AdminRole || book.SellerID != seller.ID {
				return fmt.Errorf("current user is not admin or seller of book with id %s", bookId)
			}

			return bookRepository.DeleteBook(ctx, bookId)
		},
	)
}
