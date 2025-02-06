package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetInventoryResponse struct {
	ID              uuid.UUID      `json:"id"`
	BookID          uuid.UUID      `json:"book_id"`
	Stock           int            `json:"stock"`
	LastRestockedAt time.Time      `json:"last_restocked_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty"`
}
