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
	userController, _ := wire.InitializeUserController(db)

	// Add base routes
	r.GET("/", baseController.HomePage)
	r.GET("/health", baseController.HealthCheck)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	// User routes
	users := r.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	{
		users.GET("", userController.GetUsers)
		users.GET("/me", userController.GetCurrentUser)
		users.GET(":id", userController.GetUserById)
		users.PATCH(":id", userController.UpdateUser)
		users.DELETE(":id", userController.DeleteUser)
	}

	api := r.Group("/api")
	{
		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productController.GetProducts)
			products.GET(":id", productController.GetProduct)
			products.POST("", middlewares.AuthMiddleware(), productController.CreateProduct)
			products.PATCH(":id", middlewares.AuthMiddleware(), productController.UpdateProduct)
			products.DELETE(":id", middlewares.AuthMiddleware(), productController.DeleteProduct)
		}
	}

	return r
}
