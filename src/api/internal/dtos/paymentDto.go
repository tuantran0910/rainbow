package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetPaymentResponse struct {
	ID        uuid.UUID      `json:"id"`
	Method    string         `json:"method"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListPaymentsResponse struct {
	Payments []*GetPaymentResponse `json:"payments"`
}
