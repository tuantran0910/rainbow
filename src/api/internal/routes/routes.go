package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/cmd/wire"
	"github.com/tuantran0910/rainbow/internal/controllers"
	"github.com/tuantran0910/rainbow/internal/middlewares"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Define middlewares
	r.Use(middlewares.LoggingMiddleware())

	// Define controllers
	baseController, _ := controllers.NewBaseController()
	productController, _ := wire.InitializeProductController(db)
	authController, _ := wire.InitializeAuthController(db)

	// Add base routes
	r.GET("/", baseController.HomePage)
	r.GET("/health", baseController.HealthCheck)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	api := r.Group("/api")
	{
		// Product routes
		products := api.Group("/products")
		products.Use(middlewares.AuthMiddleware())
		{
			products.GET("", productController.GetProducts)
			products.GET(":id", productController.GetProduct)
			products.POST("", productController.CreateProduct)
			products.PATCH(":id", productController.UpdateProduct)
			products.DELETE(":id", productController.DeleteProduct)
		}
	}

	return r
}
