package models

import "github.com/google/uuid"

type Seller struct {
	BaseGormModel
	UserID uuid.UUID `json:"user_id" gorm:"unique;not null"`
	Name   string    `json:"name"`
	Link   string    `json:"link" gorm:"unique;null"`
	Logo   string    `json:"logo" gorm:"unique;null"`
}
