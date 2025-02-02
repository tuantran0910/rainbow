//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	authController "github.com/tuantran0910/rainbow/internal/controllers/auth"
	productController "github.com/tuantran0910/rainbow/internal/controllers/v1/product"
	productRepository "github.com/tuantran0910/rainbow/internal/repositories/product"
	userRepository "github.com/tuantran0910/rainbow/internal/repositories/user"
	authService "github.com/tuantran0910/rainbow/internal/services/auth"
	productService "github.com/tuantran0910/rainbow/internal/services/product"
)

func InitializeProductController(db *gorm.DB) (*productController.ProductController, error) {
	wire.Build(
		productRepository.NewProductRepository,
		productService.NewProductService,
		productController.NewProductController,
	)

	return &productController.ProductController{}, nil
}

func InitializeAuthController(db *gorm.DB) (*authController.AuthController, error) {
	wire.Build(
		userRepository.NewUserRepository,
		authService.NewAuthService,
		authController.NewAuthController,
	)

	return &authController.AuthController{}, nil
}
