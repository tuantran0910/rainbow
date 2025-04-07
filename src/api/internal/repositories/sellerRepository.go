package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type ISellerRepository interface {
	WithTX(tx *gorm.DB) ISellerRepository
	GetSellers(
		ctx context.Context,
		page, limit int,
	) ([]*models.Seller, *pagination.Pagination, error)
	GetSellerById(ctx context.Context, id interface{}, isSecondary bool) (*models.Seller, error)
	GetSellerByUserId(ctx context.Context, userId uuid.UUID) (*models.Seller, error)
	CreateSeller(ctx context.Context, seller *models.Seller) error
	UpdateSeller(ctx context.Context, id interface{}, seller *models.Seller, isSecondary bool) error
	DeleteSeller(ctx context.Context, sellerId uuid.UUID) error
}

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) ISellerRepository {
	return &sellerRepository{
		db: db,
	}
}

func (sr *sellerRepository) WithTX(tx *gorm.DB) ISellerRepository {
	if tx == nil {
		return sr
	}
	return &sellerRepository{
		db: tx,
	}
}

func (sr *sellerRepository) GetSellers(
	ctx context.Context,
	page, limit int,
) ([]*models.Seller, *pagination.Pagination, error) {
	var totalSellers int64
	if err := sr.db.WithContext(ctx).Model(&models.Seller{}).Count(&totalSellers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of sellers: %w", err)
	}

	pagination := pagination.NewPagination(page, limit, int(totalSellers))

	var sellers []*models.Seller
	if err := sr.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&sellers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch sellers: %w", err)
	}
	return sellers, pagination, nil
}

func (sr *sellerRepository) GetSellerById(
	ctx context.Context,
	id interface{},
	isSecondary bool,
) (*models.Seller, error) {
	var seller models.Seller
	query := sr.db.WithContext(ctx)

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.Take(&seller).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		idStr := fmt.Sprintf("%v", id)
		return nil, fmt.Errorf("failed to get seller with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, err)
	}
	return &seller, nil
}

func (sr *sellerRepository) GetSellerByUserId(
	ctx context.Context,
	userId uuid.UUID,
) (*models.Seller, error) {
	var seller models.Seller
	if err := sr.db.WithContext(ctx).Take(&seller, "user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get seller: %w", err)
	}
	return &seller, nil
}

func (sr *sellerRepository) CreateSeller(ctx context.Context, seller *models.Seller) error {
	if err := sr.db.WithContext(ctx).Create(seller).Error; err != nil {
		return fmt.Errorf("failed to create seller: %w", err)
	}
	return nil
}

func (sr *sellerRepository) UpdateSeller(
	ctx context.Context,
	id interface{},
	seller *models.Seller,
	isSecondary bool,
) error {
	query := sr.db.WithContext(ctx).Model(&models.Seller{})

	if isSecondary {
		query = query.Where("secondary_id = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	result := query.Updates(seller)
	if result.Error != nil {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("failed to update seller with %s %s: %w",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr, result.Error)
	}

	if result.RowsAffected == 0 {
		idStr := fmt.Sprintf("%v", id)
		return fmt.Errorf("seller with %s %s not found",
			map[bool]string{true: "secondary id", false: "id"}[isSecondary],
			idStr)
	}
	return nil
}

func (sr *sellerRepository) DeleteSeller(ctx context.Context, sellerId uuid.UUID) error {
	result := sr.db.WithContext(ctx).Unscoped().Where("id = ?", sellerId).Delete(&models.Seller{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete seller: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("seller with id %s not found", sellerId)
	}
	return nil
}
