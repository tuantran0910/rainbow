package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type ISellerService interface {
	GetSellers(ctx context.Context, page, limit int) ([]*models.Seller, *pagination.Pagination, error)
	GetSellerById(ctx context.Context, sellerId uuid.UUID) (*models.Seller, error)
	CreateSeller(ctx context.Context, sellerRequest dtos.CreateSellerRequest, currentUserId uuid.UUID) error
	UpdateSeller(ctx context.Context, sellerId uuid.UUID, sellerRequest dtos.UpdateSellerRequest, currentUserId uuid.UUID) error
	DeleteSeller(ctx context.Context, sellerId uuid.UUID, currentUserId uuid.UUID) error
}

type sellerService struct {
	db               *gorm.DB
	sellerRepository repositories.ISellerRepository
	userRepository   repositories.IUserRepository
}

func NewSellerService(db *gorm.DB, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) ISellerService {
	return &sellerService{
		db:               db,
		sellerRepository: sellerRepository,
		userRepository:   userRepository,
	}
}

func (ss *sellerService) withTx(ctx context.Context, fn func(context.Context, repositories.ISellerRepository, repositories.IUserRepository) error) error {
	return ss.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sellerRepository := ss.sellerRepository.WithTX(tx)
		userRepository := ss.userRepository.WithTX(tx)
		return fn(ctx, sellerRepository, userRepository)
	})
}

func (ss *sellerService) GetSellers(ctx context.Context, page, limit int) ([]*models.Seller, *pagination.Pagination, error) {
	return ss.sellerRepository.GetSellers(ctx, page, limit)
}

func (ss *sellerService) GetSellerById(ctx context.Context, sellerId uuid.UUID) (*models.Seller, error) {
	return ss.sellerRepository.GetSellerById(ctx, sellerId)
}

func (ss *sellerService) CreateSeller(ctx context.Context, sellerRequest dtos.CreateSellerRequest, currentUserId uuid.UUID) error {
	return ss.withTx(ctx, func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
		exitingSeller, err := sellerRepository.GetSellerByUserId(ctx, currentUserId)
		if err != nil {
			if exitingSeller != nil {
				return fmt.Errorf("seller already exists")
			}
			return err
		}

		if sellerRequest.Link == "" {
			sellerRequest.Link = fmt.Sprintf("https://rainbow.tuantrann.work/sellers/%s-%s", slug.Make(sellerRequest.Name), uuid.New().String()[:8])
		}

		seller := &models.Seller{
			UserID: currentUserId,
			Link:   sellerRequest.Link,
			Name:   sellerRequest.Name,
			Logo:   sellerRequest.Logo,
		}
		return ss.sellerRepository.CreateSeller(ctx, seller)
	})
}

func (ss *sellerService) UpdateSeller(ctx context.Context, sellerId uuid.UUID, sellerRequest dtos.UpdateSellerRequest, currentUserId uuid.UUID) error {
	return ss.withTx(ctx, func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
		seller, err := sellerRepository.GetSellerById(ctx, sellerId)
		if err != nil {
			return err
		}
		if seller == nil {
			return fmt.Errorf("seller with id %s not found", sellerId)
		}

		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.ID != seller.UserID && currentUser.Role != models.AdminRole {
			return fmt.Errorf("only admin can update other seller's information")
		}

		if sellerRequest.Name != nil && *sellerRequest.Name != seller.Name {
			seller.Name = *sellerRequest.Name
			seller.Link = fmt.Sprintf("https://rainbow.tuantrann.work/sellers/%s-%s", slug.Make(*sellerRequest.Name), uuid.New().String()[:8])
		}
		if sellerRequest.Logo != nil && *sellerRequest.Logo != seller.Logo {
			seller.Logo = *sellerRequest.Logo
		}
		return sellerRepository.UpdateSeller(ctx, sellerId, seller)
	})
}

func (ss *sellerService) DeleteSeller(ctx context.Context, sellerId uuid.UUID, currentUserId uuid.UUID) error {
	return ss.withTx(ctx, func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
		seller, err := sellerRepository.GetSellerById(ctx, sellerId)
		if err != nil {
			return err
		}
		if seller == nil {
			return fmt.Errorf("seller with id %s not found", sellerId)
		}

		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil {
			return err
		}
		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.ID != seller.UserID && currentUser.Role != models.AdminRole {
			return fmt.Errorf("only admin can delete other seller")
		}
		return sellerRepository.DeleteSeller(ctx, sellerId)
	})
}
