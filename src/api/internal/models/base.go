package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseGormModel struct {
	ID        uuid.UUID      `json:"id"         gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (b *BaseGormModel) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
