package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
	Slug string `json:"slug" binding:"omitempty,min=1,max=255"`
}

type UpdateAuthorRequest struct {
	Name *string `json:"name" binding:"omitempty,min=1,max=255"`
}

type GetAuthorResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListAuthorsResponse struct {
	Authors []*GetAuthorResponse `json:"authors"`
}
