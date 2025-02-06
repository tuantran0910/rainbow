//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/tuantran0910/rainbow/internal/controllers"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/internal/services"
	"gorm.io/gorm"
)

func InitializeBookController(db *gorm.DB) (*controllers.BookController, error) {
	wire.Build(
		repositories.NewBookRepository,
		repositories.NewInventoryRepository,
		repositories.NewSellerRepository,
		repositories.NewUserRepository,
		services.NewBookService,
		controllers.NewBookController,
	)

	return &controllers.BookController{}, nil
}

func InitializeAuthController(db *gorm.DB) (*controllers.AuthController, error) {
	wire.Build(
		repositories.NewUserRepository,
		services.NewAuthService,
		controllers.NewAuthController,
	)

	return &controllers.AuthController{}, nil
}

func InitializeUserController(db *gorm.DB) (*controllers.UserController, error) {
	wire.Build(
		repositories.NewUserRepository,
		services.NewUserService,
		controllers.NewUserController,
	)

	return &controllers.UserController{}, nil
}

func InitializeSellerController(db *gorm.DB) (*controllers.SellerController, error) {
	wire.Build(
		repositories.NewUserRepository,
		repositories.NewSellerRepository,
		services.NewSellerService,
		controllers.NewSellerController,
	)

	return &controllers.SellerController{}, nil
}

func InitializeCategoryController(db *gorm.DB) (*controllers.CategoryController, error) {
	wire.Build(
		repositories.NewUserRepository,
		repositories.NewCategoryRepository,
		services.NewCategoryService,
		controllers.NewCategoryController,
	)

	return &controllers.CategoryController{}, nil
}

func InitializeAuthorController(db *gorm.DB) (*controllers.AuthorController, error) {
	wire.Build(
		repositories.NewUserRepository,
		repositories.NewAuthorRepository,
		services.NewAuthorService,
		controllers.NewAuthorController,
	)

	return &controllers.AuthorController{}, nil
}
