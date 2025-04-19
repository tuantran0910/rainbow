package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IUserPromotionRepository interface {
	WithTX(tx *gorm.DB) IUserPromotionRepository
	GetUserPromotionByUserAndPromotion(
		ctx context.Context,
		userID, promotionID uuid.UUID,
	) (*models.UserPromotion, error)
	CreateUserPromotion(ctx context.Context, userPromotion *models.UserPromotion) error
}

type userPromotionRepository struct {
	db *gorm.DB
}

func NewUserPromotionRepository(db *gorm.DB) IUserPromotionRepository {
	return &userPromotionRepository{
		db: db,
	}
}

func (upr *userPromotionRepository) WithTX(tx *gorm.DB) IUserPromotionRepository {
	if tx == nil {
		return upr
	}
	return &userPromotionRepository{
		db: tx,
	}
}

func (upr *userPromotionRepository) GetUserPromotionByUserAndPromotion(
	ctx context.Context,
	userID, promotionID uuid.UUID,
) (*models.UserPromotion, error) {
	var userPromotion models.UserPromotion
	err := upr.db.WithContext(ctx).
		Where("user_id = ? AND promotion_id = ?", userID, promotionID).
		Take(&userPromotion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user promotion: %w", err)
	}
	return &userPromotion, nil
}

func (upr *userPromotionRepository) CreateUserPromotion(
	ctx context.Context,
	userPromotion *models.UserPromotion,
) error {
	if err := upr.db.WithContext(ctx).Create(userPromotion).Error; err != nil {
		return fmt.Errorf("failed to create user promotion: %w", err)
	}
	return nil
}
