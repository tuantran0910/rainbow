package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IPromotionRepository interface {
	WithTX(tx *gorm.DB) IPromotionRepository
	GetCurrentBookPromotions(ctx context.Context, bookID uuid.UUID, currentTime time.Time) ([]*models.Promotion, error)
	CreatePromotion(ctx context.Context, promotion *models.Promotion) error
	AddBookToPromotion(ctx context.Context, promotionBook *models.PromotionBook) error
}

type promotionRepository struct {
	db *gorm.DB
}

func NewPromotionRepository(db *gorm.DB) IPromotionRepository {
	return &promotionRepository{
		db: db,
	}
}

func (pr *promotionRepository) WithTX(tx *gorm.DB) IPromotionRepository {
	if tx == nil {
		return pr
	}
	return &promotionRepository{
		db: tx,
	}
}

func (pr *promotionRepository) GetCurrentBookPromotions(ctx context.Context, bookID uuid.UUID, currentTime time.Time) ([]*models.Promotion, error) {
	var promotions []*models.Promotion
	err := pr.db.WithContext(ctx).
		Joins("JOIN promotion_books ON promotion_books.promotion_id = promotions.id").
		Where("promotion_books.book_id = ?", bookID).
		Where("promotions.start_date <= ? AND promotions.end_date >= ?", currentTime, currentTime).
		Find(&promotions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch current book promotions: %w", err)
	}
	return promotions, nil
}

func (pr *promotionRepository) CreatePromotion(ctx context.Context, promotion *models.Promotion) error {
	if err := pr.db.WithContext(ctx).Create(promotion).Error; err != nil {
		return fmt.Errorf("failed to create promotion: %w", err)
	}
	return nil
}

func (pr *promotionRepository) AddBookToPromotion(ctx context.Context, promotionBook *models.PromotionBook) error {
	if err := pr.db.WithContext(ctx).Create(promotionBook).Error; err != nil {
		return fmt.Errorf("failed to add book to promotion: %w", err)
	}
	return nil
}
