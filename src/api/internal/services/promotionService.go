package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"gorm.io/gorm"
)

type IPromotionService interface {
	CreatePromotion(ctx context.Context, promotionRequest dtos.CreatePromotionRequest, currentUserId uuid.UUID) error
}

type promotionService struct {
	db                  *gorm.DB
	promotionRepository repositories.IPromotionRepository
	userRepository      repositories.IUserRepository
}

func NewPromotionService(db *gorm.DB, promotionRepository repositories.IPromotionRepository, userRepository repositories.IUserRepository) IPromotionService {
	return &promotionService{
		db:                  db,
		promotionRepository: promotionRepository,
		userRepository:      userRepository,
	}
}

func (ps *promotionService) withTx(ctx context.Context, fn func(context.Context, repositories.IPromotionRepository, repositories.IUserRepository) error) error {
	return ps.db.Transaction(func(tx *gorm.DB) error {
		promotionRepository := ps.promotionRepository.WithTX(tx)
		userRepository := ps.userRepository.WithTX(tx)
		return fn(ctx, promotionRepository, userRepository)
	})
}

func (pc *promotionService) CreatePromotion(ctx context.Context, promotionRequest dtos.CreatePromotionRequest, currentUserId uuid.UUID) error {
	return pc.withTx(ctx, func(ctx context.Context, promotionRepository repositories.IPromotionRepository, userRepository repositories.IUserRepository) error {
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.Role != models.AdminRole {
			return fmt.Errorf("current user does not have permission to create promotion")
		}

		promotion := &models.Promotion{
			Name:          promotionRequest.Name,
			DiscountType:  promotionRequest.DiscountType,
			DiscountValue: promotionRequest.DiscountValue,
			StartDate:     promotionRequest.StartDate,
			EndDate:       promotionRequest.EndDate,
			MaxUses:       promotionRequest.MaxUses,
		}
		if err := promotionRepository.CreatePromotion(ctx, promotion); err != nil {
			return err
		}

		for _, bookId := range promotionRequest.BookIds {
			promotionBook := &models.PromotionBook{
				PromotionID: promotion.ID,
				BookID:      bookId,
			}
			if err := promotionRepository.AddBookToPromotion(ctx, promotionBook); err != nil {
				return err
			}
		}
		return nil
	})
}
