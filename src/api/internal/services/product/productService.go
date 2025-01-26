package product

import (
	"context"

	"github.com/google/uuid"
	model "github.com/tuantran0910/rainbow/internal/models/product"
	repository "github.com/tuantran0910/rainbow/internal/repositories/product"
	"github.com/tuantran0910/rainbow/pkg/pagination"
)

type IProductService interface {
	GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error)
	GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error)
	CreateProduct(ctx context.Context, productRequest model.ProductRequest) error
	UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type ProductService struct {
	productRepository repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (ps *ProductService) GetProducts(ctx context.Context, page, limit int) ([]*model.Product, *pagination.Pagination, error) {
	return ps.productRepository.GetProducts(ctx, page, limit)
}

func (ps *ProductService) GetProduct(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	return ps.productRepository.GetProduct(ctx, productId)
}

func (ps *ProductService) CreateProduct(ctx context.Context, productRequest model.ProductRequest) error {
	// Define a new product
	product := &model.Product{
		Name:  *productRequest.Name,
		Price: *productRequest.Price,
	}

	// Begin a transaction
	productRepository, err := ps.productRepository.Begin()
	if err != nil {
		return err
	}

	// Ensure the transaction is committed or rolled back properly
	defer func() {
		if r := recover(); r != nil {
			productRepository.Rollback()
			panic(r)
		} else if err != nil {
			productRepository.Rollback()
		} else {
			productRepository.Commit()
		}
	}()

	if err := productRepository.CreateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, productId uuid.UUID, productRequest model.ProductRequest) error {
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

	// Begin a transaction
	productRepository, err := ps.productRepository.Begin()
	if err != nil {
		return err
	}

	// Ensure the transaction is committed or rolled back properly
	defer func() {
		if r := recover(); r != nil {
			productRepository.Rollback()
			panic(r)
		} else if err != nil {
			productRepository.Rollback()
		} else {
			productRepository.Commit()
		}
	}()

	if err := productRepository.UpdateProduct(ctx, productId, product); err != nil {
		return err
	}

	return nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	// Begin a transaction
	productRepository, err := ps.productRepository.Begin()
	if err != nil {
		return err
	}

	// Ensure the transaction is committed or rolled back properly
	defer func() {
		if r := recover(); r != nil {
			productRepository.Rollback()
			panic(r)
		} else if err != nil {
			productRepository.Rollback()
		} else {
			productRepository.Commit()
		}
	}()

	// Delete a product
	if err := productRepository.DeleteProduct(ctx, productId); err != nil {
		return err
	}

	return nil
}
