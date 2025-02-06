package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	BaseGormModel
	BookID          uuid.UUID `json:"book_id" gorm:"unique;not null"`
	Stock           int       `json:"stock" gorm:"not null;check:stock >= 0"`
	LastRestockedAt time.Time `json:"last_restocked_at"`
}
