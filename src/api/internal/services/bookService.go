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
	GetBookById(ctx context.Context, bookId uuid.UUID) (*models.Book, error)
	CreateBook(ctx context.Context, bookRequest dtos.CreateBookRequest, currentUserId uuid.UUID) error
	UpdateBook(ctx context.Context, bookId uuid.UUID, bookRequest dtos.UpdateBookRequest, currentUserId uuid.UUID) error
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

func (ps *bookService) withTx(ctx context.Context, fn func(context.Context, repositories.IBookRepository, repositories.IInventoryRepository, repositories.ISellerRepository, repositories.IUserRepository) error) error {
	return ps.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bookRepository := ps.bookRepository.WithTX(tx)
		inventoryRepository := ps.inventoryRepository.WithTX(tx)
		sellerRepository := ps.sellerRepository.WithTX(tx)
		userRepository := ps.userRepository.WithTX(tx)
		return fn(ctx, bookRepository, inventoryRepository, sellerRepository, userRepository)
	})
}

func (ps *bookService) GetBooks(ctx context.Context, page, limit int) ([]*models.Book, *pagination.Pagination, error) {
	return ps.bookRepository.GetBooks(ctx, page, limit)
}

func (ps *bookService) GetBookById(ctx context.Context, bookId uuid.UUID) (*models.Book, error) {
	return ps.bookRepository.GetBookById(ctx, bookId)
}

func (ps *bookService) CreateBook(ctx context.Context, bookRequest dtos.CreateBookRequest, currentUserId uuid.UUID) error {
	return ps.withTx(ctx, func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
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

		if bookRequest.SellerID != seller.ID {
			return fmt.Errorf("seller id %s does not match with current user seller id %s", bookRequest.SellerID, seller.ID)
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
		}
		if err := bookRepository.CreateBook(ctx, book); err != nil {
			return err
		}

		bookStock := bookRequest.Stock
		if bookStock == 0 {
			bookStock = 100
		}
		inventory := &models.Inventory{
			BookID: book.ID,
			Stock:  bookStock,
		}
		if err := inventoryRepository.CreateInventory(ctx, inventory); err != nil {
			return err
		}

		for _, authorID := range bookRequest.AuthorIds {
			bookAuthor := &models.BookAuthor{
				BookID:   book.ID,
				AuthorID: authorID,
			}
			if err := bookRepository.AddAuthorToBook(ctx, bookAuthor); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ps *bookService) UpdateBook(ctx context.Context, bookId uuid.UUID, bookRequest dtos.UpdateBookRequest, currentUserId uuid.UUID) error {
	return ps.withTx(ctx, func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
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

		book, err := ps.bookRepository.GetBookById(ctx, bookId)
		if err != nil {
			return err
		}
		if book == nil {
			return fmt.Errorf("book with id %s not found", bookId)
		}

		if book.SellerID != seller.ID {
			return fmt.Errorf("seller id %s does not match with current user seller id %s", book.SellerID, seller.ID)
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
		if bookRequest.OriginalPrice != nil && *bookRequest.OriginalPrice != book.OriginalPrice {
			book.OriginalPrice = *bookRequest.OriginalPrice
		}
		if bookRequest.RatingAverage != nil && *bookRequest.RatingAverage != book.RatingAverage {
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
			bookInventory, err := inventoryRepository.GetInventoryByBookID(ctx, bookId)
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
			return inventoryRepository.UpdateInventory(ctx, bookInventory.ID, bookInventory)
		}
		return nil
	})
}

func (ps *bookService) DeleteBook(ctx context.Context, bookId uuid.UUID, currentUserId uuid.UUID) error {
	return ps.withTx(ctx, func(ctx context.Context, bookRepository repositories.IBookRepository, inventoryRepository repositories.IInventoryRepository, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
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

		book, err := ps.bookRepository.GetBookById(ctx, bookId)
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
	})
}
