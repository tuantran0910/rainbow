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
		// Check if user has already been a seller
		exitingSeller, err := sellerRepository.GetSellerByUserID(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if exitingSeller != nil {
			return fmt.Errorf("user has already been a seller")
		}

		// Construct a seller link
		if sellerRequest.Link == "" {
			sellerRequest.Link = fmt.Sprintf("https://rainbow.tuantrann.work/sellers/%s-%s", slug.Make(sellerRequest.Name), uuid.New().String()[:8])
		}

		// Define a new seller
		seller := &models.Seller{
			UserID: currentUserId,
			Link:   sellerRequest.Link,
			Name:   sellerRequest.Name,
			Logo:   sellerRequest.Logo,
		}

		// Create a new seller
		if err := ss.sellerRepository.CreateSeller(ctx, seller); err != nil {
			return err
		}

		return nil
	})
}

func (ss *sellerService) UpdateSeller(ctx context.Context, sellerId uuid.UUID, sellerRequest dtos.UpdateSellerRequest, currentUserId uuid.UUID) error {
	return ss.withTx(ctx, func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
		// Get seller by id
		seller, err := sellerRepository.GetSellerById(ctx, sellerId)
		if err != nil {
			return err
		}

		// Get current user by id
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.ID != seller.UserID && string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can update other seller's information")
		}

		// Update seller
		if sellerRequest.Name != nil && *sellerRequest.Name != "" {
			seller.Name = *sellerRequest.Name
			seller.Link = fmt.Sprintf("https://rainbow.tuantrann.work/sellers/%s-%s", slug.Make(*sellerRequest.Name), uuid.New().String()[:8])
		}
		if sellerRequest.Logo != nil {
			seller.Logo = *sellerRequest.Logo
		}

		// Save seller
		if err := sellerRepository.UpdateSeller(ctx, sellerId, seller); err != nil {
			return err
		}

		return nil
	})
}

func (ss *sellerService) DeleteSeller(ctx context.Context, sellerId uuid.UUID, currentUserId uuid.UUID) error {
	return ss.withTx(ctx, func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
		// Get seller by id
		seller, err := sellerRepository.GetSellerById(ctx, sellerId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if seller == nil {
			return fmt.Errorf("seller not found")
		}

		// Get current user by id
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.ID != seller.UserID && string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can delete other seller")
		}

		// Delete a seller
		return sellerRepository.DeleteSeller(ctx, sellerId)
	})
}
