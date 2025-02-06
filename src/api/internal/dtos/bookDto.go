package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateBookRequest struct {
	SecondaryID   string      `json:"secondary_id" binding:"omitempty"`
	CategoryID    uuid.UUID   `json:"category_id" binding:"required"`
	SellerID      uuid.UUID   `json:"seller_id" binding:"required"`
	Name          string      `json:"name" binding:"required,min=1,max=255"`
	Description   string      `json:"description" binding:"omitempty"`
	Price         float64     `json:"price" binding:"required,gt=0"`
	OriginalPrice float64     `json:"original_price" binding:"omitempty,gt=0"`
	RatingAverage float64     `json:"rating_average" binding:"omitempty,gt=0"`
	ReviewCount   int         `json:"review_count" binding:"omitempty,gt=0"`
	PageCount     int         `json:"page_count" binding:"omitempty,gt=0"`
	Stock         int         `json:"stock" binding:"omitempty,gt=0"`
	AuthorIds     []uuid.UUID `json:"author_ids" binding:"required"`
}

type UpdateBookRequest struct {
	CategoryID    *uuid.UUID `json:"category_id" binding:"omitempty"`
	Name          *string    `json:"name" binding:"omitempty,min=1,max=255"`
	Description   *string    `json:"description" binding:"omitempty"`
	Price         *float64   `json:"price" binding:"omitempty,gt=0"`
	OriginalPrice *float64   `json:"original_price" binding:"omitempty,gt=0"`
	RatingAverage *float64   `json:"rating_average" binding:"omitempty,gt=0"`
	ReviewCount   *int       `json:"review_count" binding:"omitempty,gt=0"`
	PageCount     *int       `json:"page_count" binding:"omitempty,gt=0"`
	Stock         *int       `json:"stock" binding:"omitempty,gt=0"`
}

type GetBookResponse struct {
	ID            uuid.UUID             `json:"id"`
	SecondaryID   string                `json:"secondary_id"`
	CategoryID    uuid.UUID             `json:"category_id"`
	SellerID      uuid.UUID             `json:"seller_id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	Price         float64               `json:"price"`
	OriginalPrice float64               `json:"original_price"`
	RatingAverage float64               `json:"rating_average"`
	ReviewCount   int                   `json:"review_count"`
	PageCount     int                   `json:"page_count"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `json:"deleted_at,omitempty"`
	Stock         *GetInventoryResponse `json:"stock,omitempty"`
	Authors       []*GetAuthorResponse  `json:"authors"`
}

type ListBooksResponse struct {
	Books []*GetBookResponse `json:"books"`
}
