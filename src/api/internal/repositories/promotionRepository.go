package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IPromotionRepository interface {
	WithTX(tx *gorm.DB) IPromotionRepository
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
