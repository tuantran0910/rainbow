package models

import (
	"time"

	"github.com/google/uuid"
)

type DiscountType string

const (
	Percentage DiscountType = "PERCENTAGE"
	Fixed      DiscountType = "FIXED"
)

type Promotion struct {
	BaseGormModel
	Name          string       `json:"name" gorm:"not null"`
	DiscountType  DiscountType `json:"discount_type" gorm:"type:enum('PERCENTAGE', 'FIXED')"`
	DiscountValue float64      `json:"discount_value" gorm:"not null;check:discount_value >= 0"`
	StartDate     time.Time    `json:"start_date" gorm:"not null"`
	EndDate       time.Time    `json:"end_date" gorm:"not null;check:end_date >= start_date"`
	MaxUses       int          `json:"max_uses" gorm:"not null;check:max_uses >= 0"`
	UsedCount     int          `json:"used_count" gorm:"not null;check:used_count >= 0;default:0"`
	Books         []Book       `json:"books" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:promotion_books"`
}

type PromotionBook struct {
	BaseGormModel
	BookID      uuid.UUID `json:"book_id" gorm:"not null"`
	PromotionID uuid.UUID `json:"promotion_id" gorm:"not null"`
}
