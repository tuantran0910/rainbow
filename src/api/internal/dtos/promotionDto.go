package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type CreatePromotionRequest struct {
	Name          string              `json:"name"           binding:"required"`
	DiscountType  models.DiscountType `json:"discount_type"  binding:"required,oneof=PERCENTAGE FIXED"`
	DiscountValue float64             `json:"discount_value" binding:"required"`
	StartDate     time.Time           `json:"start_date"     binding:"required"`
	EndDate       time.Time           `json:"end_date"       binding:"required"`
	MaxUses       int                 `json:"max_uses"       binding:"required"`
}

type GetPromotionResponse struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	DiscountType  models.DiscountType `json:"discount_type"`
	DiscountValue float64             `json:"discount_value"`
	StartDate     time.Time           `json:"start_date"`
	EndDate       time.Time           `json:"end_date"`
	MaxUses       int                 `json:"max_uses"`
	UsedCount     int                 `json:"used_count"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `json:"deleted_at,omitempty"`
}

type ListPromotionsResponse struct {
	Promotions []*GetPromotionResponse `json:"promotions"`
}
