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
	CreateSeller(
		ctx context.Context,
		sellerReq dtos.CreateSellerRequest,
	) (*models.Seller, error)
	UpdateSeller(
		ctx context.Context,
		id interface{},
		sellerReq dtos.UpdateSellerRequest,
		isSecondary bool,
	) (*models.Seller, error)
	DeleteSeller(ctx context.Context, sellerId uuid.UUID) error
}

type sellerService struct {
	db               *gorm.DB
	sellerRepository repositories.ISellerRepository
}

func NewSellerService(
	db *gorm.DB,
	sellerRepository repositories.ISellerRepository,
) ISellerService {
	return &sellerService{
		db:               db,
		sellerRepository: sellerRepository,
	}
}

func (ss *sellerService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.ISellerRepository) error,
) error {
	return ss.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sellerRepository := ss.sellerRepository.WithTX(tx)
		return fn(ctx, sellerRepository)
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

func (ss *sellerService) CreateSeller(
	ctx context.Context,
	sellerReq dtos.CreateSellerRequest,
) (*models.Seller, error) {
	var createdSeller *models.Seller
	err := ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository) error {
			if sellerReq.Link == "" {
				sellerReq.Link = fmt.Sprintf(
					"https://rainbow.tuantrann.work/sellers/%s-%s",
					slug.Make(sellerReq.Name),
					uuid.New().String()[:8],
				)
			}

			seller := &models.Seller{
				SecondaryID: sellerReq.SecondaryID,
				Name:        sellerReq.Name,
				Link:        sellerReq.Link,
				Logo:        sellerReq.Logo,
			}
			if err := sellerRepository.CreateSeller(ctx, seller); err != nil {
				return err
			}

			// Get the created seller with all fields populated
			var getErr error
			createdSeller, getErr = sellerRepository.GetSellerById(ctx, seller.ID, false)
			return getErr
		},
	)
	return createdSeller, err
}

func (ss *sellerService) UpdateSeller(
	ctx context.Context,
	id interface{},
	sellerReq dtos.UpdateSellerRequest,
	isSecondary bool,
) (*models.Seller, error) {
	var updatedSeller *models.Seller
	err := ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository) error {
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

			var sellerToUpdate models.Seller
			if sellerReq.SecondaryID != nil {
				sellerToUpdate.SecondaryID = *sellerReq.SecondaryID
			}
			if sellerReq.Name != nil && *sellerReq.Name != seller.Name {
				sellerToUpdate.Name = *sellerReq.Name
				sellerToUpdate.Link = fmt.Sprintf(
					"https://rainbow.tuantrann.work/sellers/%s-%s",
					slug.Make(*sellerReq.Name),
					uuid.New().String()[:8],
				)
			}
			if sellerReq.Logo != nil && *sellerReq.Logo != seller.Logo {
				sellerToUpdate.Logo = *sellerReq.Logo
			}
			if err := sellerRepository.UpdateSeller(ctx, id, &sellerToUpdate, isSecondary); err != nil {
				return err
			}

			// Get the updated seller
			var getErr error
			updatedSeller, getErr = sellerRepository.GetSellerById(ctx, id, isSecondary)
			return getErr
		},
	)
	return updatedSeller, err
}

func (ss *sellerService) DeleteSeller(
	ctx context.Context,
	sellerId uuid.UUID,
) error {
	return ss.withTX(
		ctx,
		func(ctx context.Context, sellerRepository repositories.ISellerRepository) error {
			seller, err := sellerRepository.GetSellerById(ctx, sellerId, false)
			if err != nil {
				return err
			}
			if seller == nil {
				return fmt.Errorf("seller with id %s not found", sellerId)
			}
			return sellerRepository.DeleteSeller(ctx, sellerId)
		},
	)
}
