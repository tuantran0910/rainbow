//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	controller "github.com/tuantran0910/rainbow/internal/controllers/v1/product"
	repository "github.com/tuantran0910/rainbow/internal/repositories/product"
	service "github.com/tuantran0910/rainbow/internal/services/product"
)

func InitializeProductController(db *gorm.DB) (*controller.ProductController, error) {
	wire.Build(
		repository.NewProductRepository,
		service.NewProductService,
		controller.NewProductController,
	)

	return &controller.ProductController{}, nil
}
