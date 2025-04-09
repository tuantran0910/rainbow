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

type IBookRepository interface {
	WithTX(tx *gorm.DB) IBookRepository
	GetBooks(ctx context.Context, page, limit int) ([]*models.Book, *pagination.Pagination, error)
	GetBookById(ctx context.Context, id interface{}, isSecondary bool) (*models.Book, error)
	CreateBook(ctx context.Context, book *models.Book) error
	AddAuthorToBook(ctx context.Context, bookAuthor *models.BookAuthor) error
	ClearBookAuthors(ctx context.Context, bookId uuid.UUID) error
	UpdateBook(ctx context.Context, bookId uuid.UUID, book *models.Book) error
	DeleteBook(ctx context.Context, bookId uuid.UUID) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) IBookRepository {
	return &bookRepository{
		db: db,
	}
}

func (pr *bookRepository) WithTX(tx *gorm.DB) IBookRepository {
	if tx == nil {
		return pr
	}

	return &bookRepository{
		db: tx,
	}
}

func (pr *bookRepository) GetBooks(
	ctx context.Context,
	page, limit int,
) ([]*models.Book, *pagination.Pagination, error) {
	var totalBooks int64
	if err := pr.db.WithContext(ctx).Model(&models.Book{}).Count(&totalBooks).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of books: %w", err)
	}

	pagination := pagination.NewPagination(page, limit, int(totalBooks))

	var books []*models.Book
	if err := pr.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch books: %w", err)
	}
	return books, pagination, nil
}

func (pr *bookRepository) GetBookById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Book, error) {
	var book models.Book
	query := pr.db.WithContext(ctx).Preload("Inventory").Preload("Authors")

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.Take(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		idStr := fmt.Sprintf("%v", id)
		return nil, fmt.Errorf("failed to fetch book with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, err)
	}
	return &book, nil
}

func (pr *bookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	if err := pr.db.WithContext(ctx).Create(book).Error; err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	return nil
}

func (pr *bookRepository) AddAuthorToBook(
	ctx context.Context,
	bookAuthor *models.BookAuthor,
) error {
	if err := pr.db.WithContext(ctx).Create(bookAuthor).Error; err != nil {
		return fmt.Errorf(
			"failed to add author with id %s to book with id %s: %w",
			bookAuthor.AuthorID,
			bookAuthor.BookID,
			err,
		)
	}
	return nil
}

func (pr *bookRepository) ClearBookAuthors(
	ctx context.Context,
	bookId uuid.UUID,
) error {
	if err := pr.db.WithContext(ctx).Where("book_id = ?", bookId).Delete(&models.BookAuthor{}).Error; err != nil {
		return fmt.Errorf("failed to clear authors for book with id %s: %w", bookId, err)
	}
	return nil
}

func (pr *bookRepository) UpdateBook(
	ctx context.Context,
	bookId uuid.UUID,
	book *models.Book,
) error {
	result := pr.db.WithContext(ctx).
		Model(&models.Book{}).
		Where("id = ?", bookId).
		Select("SecondaryID", "CategoryID", "Name", "Description", "Price", "OriginalPrice", "RatingAverage", "ReviewCount", "PageCount", "SoldCount").
		Updates(book)
	if result.Error != nil {
		return fmt.Errorf("failed to update book with id %s: %w", bookId, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("book with id %s not found", bookId)
	}
	return nil
}

func (pr *bookRepository) DeleteBook(ctx context.Context, bookId uuid.UUID) error {
	result := pr.db.WithContext(ctx).Unscoped().Where("id = ?", bookId).Delete(&models.Book{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete book with id %s: %w", bookId, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("book with id %s not found", bookId)
	}
	return nil
}
