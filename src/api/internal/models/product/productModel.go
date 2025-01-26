package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type Product struct {
	models.BaseGormModel
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductRequest struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
}

type ProductResponse struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Price     float64        `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
