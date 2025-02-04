package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateSellerRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
	Link string `json:"link"`
	Logo string `json:"logo"`
}

type UpdateSellerRequest struct {
	Name *string `json:"name" binding:"omitempty,min=1,max=255"`
	Logo *string `json:"logo" binding:"omitempty"`
}

type GetSellerResponse struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	Name      string         `json:"name"`
	Link      string         `json:"link"`
	Logo      string         `json:"logo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListSellersResponse struct {
	Sellers []*GetSellerResponse `json:"sellers"`
}
