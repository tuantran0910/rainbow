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
	GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error)
	GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error)
	CreateProduct(ctx context.Context, product *model.Product) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, product *model.Product) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (pr *productRepository) GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error) {
	// Get total number of products
	var totalProducts int64
	if err := pr.DB.WithContext(ctx).Model(&model.Product{}).Count(&totalProducts).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of products: %w", err)
	}

	// Define pagination
	pagination := pagination.NewPagination(page, limit, int(totalProducts))

	var products []*model.Product
	if err := pr.DB.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, pagination, nil
}

func (pr *productRepository) GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	var product model.Product
	if err := pr.DB.WithContext(ctx).First(&product, "id = ?", productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	return &product, nil
}

func (pr *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	if err := pr.DB.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, productId uuid.UUID, product *model.Product) error {
	if err := pr.DB.WithContext(ctx).Model(&model.Product{}).Where("id = ?", productId).Updates(product).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	if err := pr.DB.WithContext(ctx).Delete(&model.Product{}, "id = ?", productId).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
