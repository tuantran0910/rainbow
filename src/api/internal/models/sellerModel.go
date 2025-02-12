package models

import "github.com/google/uuid"

type Seller struct {
	BaseGormModel
	UserID uuid.UUID `json:"user_id" gorm:"unique;not null"`
	Name   string    `json:"name"    gorm:"not null;size:255;check:length(name) > 0"`
	Link   string    `json:"link"    gorm:"unique;not null"`
	Logo   string    `json:"logo"    gorm:"unique;null"`
	Books  []Book    `               gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
