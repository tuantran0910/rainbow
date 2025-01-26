package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/internal/controllers"
	"github.com/tuantran0910/rainbow/internal/controllers/v1/product"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Define controllers
	baseController := controllers.NewBaseController()
	productController := product.NewProductController(db)

	// Add base routes
	r.GET("/", baseController.HomePage)
	r.GET("/health", baseController.HealthCheck)

	// Add routes
	v1 := r.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
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
