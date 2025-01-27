package product

import (
	"context"

	"github.com/google/uuid"
	model "github.com/tuantran0910/rainbow/internal/models/product"
	repository "github.com/tuantran0910/rainbow/internal/repositories/product"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IProductService interface {
	GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error)
	GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error)
	CreateProduct(ctx context.Context, productRequest model.ProductRequest) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type productService struct {
	db                *gorm.DB
	productRepository repository.IProductRepository
}

func NewProductService(db *gorm.DB, productRepository repository.IProductRepository) IProductService {
	return &productService{
		db:                db,
		productRepository: productRepository,
	}
}

func (ps *productService) withTx(ctx context.Context, fn func(context.Context, repository.IProductRepository) error) error {
	return ps.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		productRepository := ps.productRepository.WithTX(tx)
		return fn(ctx, productRepository)
	})
}

func (ps *productService) GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error) {
	return ps.productRepository.GetProducts(ctx, page, limit)
}

func (ps *productService) GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	return ps.productRepository.GetProduct(ctx, productId)
}

func (ps *productService) CreateProduct(ctx context.Context, productRequest model.ProductRequest) error {
	return ps.withTx(ctx, func(ctx context.Context, productRepository repository.IProductRepository) error {
		// Define a new product
		product := &model.Product{
			Name:  *productRequest.Name,
			Price: *productRequest.Price,
		}

		// Create a new product
		if err := productRepository.CreateProduct(ctx, product); err != nil {
			return err
		}

		return nil
	})
}

func (ps *productService) UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error {
	return ps.withTx(ctx, func(ctx context.Context, productRepository repository.IProductRepository) error {
		// Define an update product
		product, err := ps.productRepository.GetProduct(ctx, productId)
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
	return ps.withTx(ctx, func(ctx context.Context, productRepository repository.IProductRepository) error {
		// Delete a product
		if err := productRepository.DeleteProduct(ctx, productId); err != nil {
			return err
		}

		return nil
	})
}
