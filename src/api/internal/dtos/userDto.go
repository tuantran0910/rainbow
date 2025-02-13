package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type UpdateUserRequest struct {
	Email       *string `json:"email"        binding:"omitempty,email"`
	FirstName   *string `json:"first_name"   binding:"omitempty,min=1,max=255"`
	LastName    *string `json:"last_name"    binding:"omitempty,min=1,max=255"`
	IsActive    *bool   `json:"is_active"    binding:"omitempty"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,len=10"`
	Role        *string `json:"role"         binding:"omitempty,oneof=ADMIN USER"`
}

type GetUserResponse struct {
	ID          uuid.UUID      `json:"id"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	LastLogin   time.Time      `json:"last_login"`
	IsActive    bool           `json:"is_active"`
	PhoneNumber string         `json:"phone_number"`
	Role        models.Role    `json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListUsersResponse struct {
	Users []*GetUserResponse `json:"users"`
}
