package repositories

import (
	"context"
	"fmt"

	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	GetPayments(ctx context.Context) ([]*models.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (pr *paymentRepository) GetPayments(ctx context.Context) ([]*models.Payment, error) {
	var totalPayments int64
	if err := pr.db.WithContext(ctx).Model(&models.Payment{}).Count(&totalPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch total number of payments: %w", err)
	}

	var payments []*models.Payment
	if err := pr.db.WithContext(ctx).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}
	return payments, nil
}
