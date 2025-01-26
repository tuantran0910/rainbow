package product

import (
	"context"

	"github.com/google/uuid"
	model "github.com/tuantran0910/rainbow/internal/models/product"
	repository "github.com/tuantran0910/rainbow/internal/repositories/product"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"github.com/tuantran0910/rainbow/pkg/utils/transaction"
	"gorm.io/gorm"
)

type IProductService interface {
	GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error)
	GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error)
	CreateProduct(ctx context.Context, productRequest model.ProductRequest) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) IProductService {
	return &ProductService{
		db: db,
	}
}

func (ps *ProductService) GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error) {
	return repository.NewProductRepository(ps.db).GetProducts(ctx, page, limit)
}

func (ps *ProductService) GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	return repository.NewProductRepository(ps.db).GetProduct(ctx, productId)
}

func (ps *ProductService) CreateProduct(ctx context.Context, productRequest model.ProductRequest) error {
	return transaction.WithTransaction(ps.db, func(tx *gorm.DB) error {
		// Define a new product
		product := &model.Product{
			Name:  *productRequest.Name,
			Price: *productRequest.Price,
		}

		if err := repository.NewProductRepository(tx).CreateProduct(ctx, product); err != nil {
			return err
		}

		return nil
	})
}

func (ps *ProductService) UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error {
	return transaction.WithTransaction(ps.db, func(tx *gorm.DB) error {
		// Define an update product
		product, err := repository.NewProductRepository(tx).GetProduct(ctx, productId)
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

		if err := repository.NewProductRepository(tx).UpdateProduct(ctx, productId, product); err != nil {
			return err
		}

		return nil
	})
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	return transaction.WithTransaction(ps.db, func(tx *gorm.DB) error {
		if err := repository.NewProductRepository(tx).DeleteProduct(ctx, productId); err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}
