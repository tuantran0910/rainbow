package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IPromotionRepository interface {
	WithTX(tx *gorm.DB) IPromotionRepository
	GetAllPromotions(ctx context.Context, userID *uuid.UUID) ([]*models.Promotion, error)
	GetPromotionById(ctx context.Context, promotionID uuid.UUID) (*models.Promotion, error)
	CreatePromotion(ctx context.Context, promotion *models.Promotion) error
	UpdatePromotion(ctx context.Context, promotionId uuid.UUID, promotion *models.Promotion) error
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

func (pr *promotionRepository) GetAllPromotions(
	ctx context.Context,
	userID *uuid.UUID,
) ([]*models.Promotion, error) {
	var promotions []*models.Promotion
	now := time.Now()

	// Start with base query for active promotions
	query := pr.db.WithContext(ctx).Where("start_date <= ? AND end_date >= ?", now, now)

	// If userID is provided, filter out promotions already used by this user
	if userID != nil {
		query = query.Where("used_count < max_uses").
			Joins("LEFT JOIN user_promotions up ON promotions.id = up.promotion_id AND up.user_id = ?", *userID).
			Where("up.id IS NULL")
	}

	err := query.Find(&promotions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch promotions: %w", err)
	}
	return promotions, nil
}

func (pr *promotionRepository) GetPromotionById(
	ctx context.Context,
	promotionID uuid.UUID,
) (*models.Promotion, error) {
	var promotion models.Promotion
	err := pr.db.WithContext(ctx).Take(&promotion, "id = ?", promotionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch promotion with id %s: %w", promotionID, err)
	}
	return &promotion, nil
}

func (pr *promotionRepository) CreatePromotion(
	ctx context.Context,
	promotion *models.Promotion,
) error {
	if err := pr.db.WithContext(ctx).Create(promotion).Error; err != nil {
		return fmt.Errorf("failed to create promotion: %w", err)
	}
	return nil
}

func (pr *promotionRepository) UpdatePromotion(
	ctx context.Context, promotionID uuid.UUID, promotion *models.Promotion,
) error {
	result := pr.db.WithContext(ctx).
		Model(&models.Promotion{}).
		Where("id = ?", promotionID).
		Select("MaxUses", "UsedCount").
		Updates(promotion)

	if result.Error != nil {
		return fmt.Errorf("failed to update promotion with id %s: %w", promotionID, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("promotion with id %s does not exist", promotionID)
	}

	return nil
}
