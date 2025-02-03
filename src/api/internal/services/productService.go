package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IProductService interface {
	GetProducts(ctx context.Context, page, limit int) ([]*models.Product, *pagination.Pagination, error)
	GetProductByID(ctx context.Context, productId uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, productRequest dtos.CreateProductRequest) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest dtos.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type productService struct {
	db                *gorm.DB
	productRepository repositories.IProductRepository
}

func NewProductService(db *gorm.DB, productRepository repositories.IProductRepository) IProductService {
	return &productService{
		db:                db,
		productRepository: productRepository,
	}
}

func (ps *productService) withTx(ctx context.Context, fn func(context.Context, repositories.IProductRepository) error) error {
	return ps.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		productRepository := ps.productRepository.WithTX(tx)
		return fn(ctx, productRepository)
	})
}

func (ps *productService) GetProducts(ctx context.Context, page, limit int) ([]*models.Product, *pagination.Pagination, error) {
	return ps.productRepository.GetProducts(ctx, page, limit)
}

func (ps *productService) GetProductByID(ctx context.Context, productId uuid.UUID) (*models.Product, error) {
	return ps.productRepository.GetProductByID(ctx, productId)
}

func (ps *productService) CreateProduct(ctx context.Context, productRequest dtos.CreateProductRequest) error {
	return ps.withTx(ctx, func(ctx context.Context, productRepository repositories.IProductRepository) error {
		// Define a new product
		product := &models.Product{
			Name:  productRequest.Name,
			Price: productRequest.Price,
		}

		// Create a new product
		if err := productRepository.CreateProduct(ctx, product); err != nil {
			return err
		}

		return nil
	})
}

func (ps *productService) UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest dtos.UpdateProductRequest) error {
	return ps.withTx(ctx, func(ctx context.Context, productRepository repositories.IProductRepository) error {
		// Define an update product
		product, err := ps.productRepository.GetProductByID(ctx, productId)
		if err != nil {
			return err
		}

		// Apply the updates
		if productRequest.Name != nil {
			product.Name = *productRequest.Name
		}
		if productRequest.Price != nil {
			product.Price = *productRequest.Price
		}

		if err := productRepository.UpdateProduct(ctx, productId, product); err != nil {
			return err
		}

		return nil
	})
}

func (ps *productService) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	return ps.withTx(ctx, func(ctx context.Context, productRepository repositories.IProductRepository) error {
		// Delete a product
		if err := productRepository.DeleteProduct(ctx, productId); err != nil {
			return err
		}

		return nil
	})
}
