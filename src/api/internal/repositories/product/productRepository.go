package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	model "github.com/tuantran0910/rainbow/internal/models/product"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IProductRepository interface {
	WithTX(tx *gorm.DB) IProductRepository
	GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error)
	GetProductByID(ctx context.Context, productId uuid.UUID) (*model.Product, error)
	CreateProduct(ctx context.Context, product *model.Product) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, product *model.Product) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}

func (pr *productRepository) WithTX(tx *gorm.DB) IProductRepository {
	if tx == nil {
		return pr
	}
	return &productRepository{
		db: tx,
	}
}

func (pr *productRepository) GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error) {
	// Get total number of products
	var totalProducts int64
	if err := pr.db.WithContext(ctx).Model(&model.Product{}).Count(&totalProducts).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of products: %w", err)
	}

	// Define pagination
	pagination := pagination.NewPagination(page, limit, int(totalProducts))

	var products []*model.Product
	if err := pr.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, pagination, nil
}

func (pr *productRepository) GetProductByID(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	var product model.Product
	if err := pr.db.WithContext(ctx).Take(&product, "id = ?", productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	return &product, nil
}

func (pr *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	if err := pr.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, productId uuid.UUID, product *model.Product) error {
	result := pr.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productId).Updates(product)
	if result.Error != nil {
		return fmt.Errorf("failed to update product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", productId)
	}

	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	result := pr.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", productId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", productId)
	}

	return nil
}
