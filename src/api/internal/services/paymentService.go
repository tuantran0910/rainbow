package services

import (
	"context"

	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
)

type IPaymentService interface {
	GetPayments(ctx context.Context) ([]*models.Payment, error)
}

type paymentService struct {
	paymentRepository repositories.IPaymentRepository
}

func NewPaymentService(paymentRepository repositories.IPaymentRepository) IPaymentService {
	return &paymentService{
		paymentRepository: paymentRepository,
	}
}

func (ps *paymentService) GetPayments(ctx context.Context) ([]*models.Payment, error) {
	return ps.paymentRepository.GetPayments(ctx)
}
