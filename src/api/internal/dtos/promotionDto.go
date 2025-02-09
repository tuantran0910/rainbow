package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
)

type CreatePromotionRequest struct {
	Name          string              `json:"name" binding:"required"`
	DiscountType  models.DiscountType `json:"discount_type" binding:"required,oneof=PERCENTAGE FIXED"`
	DiscountValue float64             `json:"discount_value" binding:"required"`
	StartDate     time.Time           `json:"start_date" binding:"required"`
	EndDate       time.Time           `json:"end_date" binding:"required"`
	MaxUses       int                 `json:"max_uses" binding:"required"`
	BookIds       []uuid.UUID         `json:"book_ids" binding:"required"`
}
