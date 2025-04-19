package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPromotion struct {
	BaseGormModel
	UserID      uuid.UUID `json:"user_id"      gorm:"not null"`
	PromotionID uuid.UUID `json:"promotion_id" gorm:"not null"`
	UsedAt      time.Time `json:"used_at"      gorm:"not null;default:CURRENT_TIMESTAMP"`
	User        User      `json:"user"         gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Promotion   Promotion `json:"promotion"    gorm:"foreignKey:PromotionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
