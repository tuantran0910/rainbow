package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=3,max=255"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name  *string  `json:"name" binding:"omitempty,min=3,max=255"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}

type GetProductResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListProductsResponse struct {
	Products []*GetProductResponse `json:"products"`
}
