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
	GetSellers(
		ctx context.Context,
		page, limit int,
	) ([]*models.Seller, *pagination.Pagination, error)
	GetSellerById(ctx context.Context, id interface{}, isSecondary bool) (*models.Seller, error)
	GetSellerByUserId(ctx context.Context, userId uuid.UUID) (*models.Seller, error)
	CreateSeller(ctx context.Context, sellerReq dtos.CreateSellerRequest, userId uuid.UUID) error
	UpdateSeller(
		ctx context.Context,
		id interface{},
		sellerReq dtos.UpdateSellerRequest,
		userId uuid.UUID,
		isSecondary bool,
	) error
	DeleteSeller(ctx context.Context, sellerId uuid.UUID, userId uuid.UUID) error
}

type sellerService struct {
	db               *gorm.DB
	sellerRepository repositories.ISellerRepository
	userRepository   repositories.IUserRepository
}

func NewSellerService(
	db *gorm.DB,
	sellerRepository repositories.ISellerRepository,
	userRepository repositories.IUserRepository,
) ISellerService {
	return &sellerService{
		db:               db,
		sellerRepository: sellerRepository,
		userRepository:   userRepository,
	}
}

func (ss *sellerService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.ISellerRepository, repositories.IUserRepository) error,
) error {
	return ss.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sellerRepository := ss.sellerRepository.WithTX(tx)
		userRepository := ss.userRepository.WithTX(tx)
		return fn(ctx, sellerRepository, userRepository)
	})
}

func (ss *sellerService) GetSellers(
	ctx context.Context,
	page, limit int,
) ([]*models.Seller, *pagination.Pagination, error) {
	return ss.sellerRepository.GetSellers(ctx, page, limit)
}

func (ss *sellerService) GetSellerById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Seller, error) {
	return ss.sellerRepository.GetSellerById(ctx, id, isSecondary)
}

func (ss *sellerService) GetSellerByUserId(
	ctx context.Context,
	userId uuid.UUID,
) (*models.Seller, error) {
	return ss.sellerRepository.GetSellerByUserId(ctx, userId)
}

func (ss *sellerService) CreateSeller(
	ctx context.Context,
	sellerReq dtos.CreateSellerRequest,
	userId uuid.UUID,
) error {
	return ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			exitingSeller, err := sellerRepository.GetSellerByUserId(ctx, userId)
			if err != nil {
				if exitingSeller != nil {
					return fmt.Errorf("seller already exists")
				}
				return err
			}

			if sellerReq.Link == "" {
				sellerReq.Link = fmt.Sprintf(
					"https://rainbow.tuantrann.work/sellers/%s-%s",
					slug.Make(sellerReq.Name),
					uuid.New().String()[:8],
				)
			}

			seller := &models.Seller{
				UserID: userId,
				Name:   sellerReq.Name,
				Link:   sellerReq.Link,
				Logo:   sellerReq.Logo,
			}
			return sellerRepository.CreateSeller(ctx, seller)
		},
	)
}

func (ss *sellerService) UpdateSeller(
	ctx context.Context,
	id interface{},
	sellerReq dtos.UpdateSellerRequest,
	userId uuid.UUID,
	isSecondary bool,
) error {
	return ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			seller, err := sellerRepository.GetSellerById(ctx, id, isSecondary)
			if err != nil {
				return err
			}
			if seller == nil {
				idStr := fmt.Sprintf("%v", id)
				return fmt.Errorf("seller with %s %s not found",
					map[bool]string{true: "secondary id", false: "id"}[isSecondary],
					idStr)
			}

			currentUser, err := userRepository.GetUserById(ctx, userId)
			if err != nil {
				return err
			}
			if currentUser == nil {
				return fmt.Errorf("current user not found")
			}

			if currentUser.ID != seller.UserID && currentUser.Role != models.AdminRole {
				return fmt.Errorf("only admin can update other seller's information")
			}

			var updatedSeller models.Seller
			if sellerReq.Name != nil && *sellerReq.Name != seller.Name {
				updatedSeller.Name = *sellerReq.Name
				updatedSeller.Link = fmt.Sprintf(
					"https://rainbow.tuantrann.work/sellers/%s-%s",
					slug.Make(*sellerReq.Name),
					uuid.New().String()[:8],
				)
			}
			if sellerReq.Logo != nil && *sellerReq.Logo != seller.Logo {
				updatedSeller.Logo = *sellerReq.Logo
			}
			return sellerRepository.UpdateSeller(ctx, id, &updatedSeller, isSecondary)
		},
	)
}

func (ss *sellerService) DeleteSeller(
	ctx context.Context,
	sellerId uuid.UUID,
	userId uuid.UUID,
) error {
	return ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository, userRepository repositories.IUserRepository) error {
			seller, err := sellerRepository.GetSellerById(ctx, sellerId, false)
			if err != nil {
				return err
			}
			if seller == nil {
				return fmt.Errorf("seller with id %s not found", sellerId)
			}

			currentUser, err := userRepository.GetUserById(ctx, userId)
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
		},
	)
}
