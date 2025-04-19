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
	bookController, _ := wire.InitializeBookController(db)
	authController, _ := wire.InitializeAuthController(db)
	userController, _ := wire.InitializeUserController(db)
	sellerController, _ := wire.InitializeSellerController(db)
	categoryController, _ := wire.InitializeCategoryController(db)
	authorController, _ := wire.InitializeAuthorController(db)
	paymentController, _ := wire.InitializePaymentController(db)
	promotionController, _ := wire.InitializePromotionController(db)
	orderController, _ := wire.InitializeOrderController(db)

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
		// Book routes
		books := api.Group("/books")
		{
			books.GET("", bookController.GetBooks)
			books.GET("/total", bookController.CountBooks)
			books.GET(":id", bookController.GetBookById)
			books.POST("", middlewares.AuthMiddleware(), bookController.CreateBook)
			books.PATCH(":id", middlewares.AuthMiddleware(), bookController.UpdateBook)
			books.DELETE(":id", middlewares.AuthMiddleware(), bookController.DeleteBook)
		}

		// Seller routes
		sellers := api.Group("/sellers")
		{
			sellers.GET("", sellerController.GetSellers)
			sellers.GET(":id", sellerController.GetSellerById)
			sellers.POST("", middlewares.AuthMiddleware(), sellerController.CreateSeller)
			sellers.PATCH(":id", middlewares.AuthMiddleware(), sellerController.UpdateSeller)
			sellers.DELETE(":id", middlewares.AuthMiddleware(), sellerController.DeleteSeller)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryController.GetCategories)
			categories.GET(":id", categoryController.GetCategoryById)
			categories.GET("/slugs/:slug", categoryController.GetCategoryBySlug)
			categories.POST("", middlewares.AuthMiddleware(), categoryController.CreateCategory)
			categories.PATCH(":id", middlewares.AuthMiddleware(), categoryController.UpdateCategory)
			categories.DELETE(
				":id",
				middlewares.AuthMiddleware(),
				categoryController.DeleteCategory,
			)
		}

		// Author routes
		authors := api.Group("/authors")
		{
			authors.GET("", authorController.GetAuthors)
			authors.GET(":id", authorController.GetAuthorById)
			authors.GET("/slugs/:slug", authorController.GetAuthorBySlug)
			authors.POST("", middlewares.AuthMiddleware(), authorController.CreateAuthor)
			authors.PATCH(":id", middlewares.AuthMiddleware(), authorController.UpdateAuthor)
			authors.DELETE(":id", middlewares.AuthMiddleware(), authorController.DeleteAuthor)
		}

		// Payment routes
		payments := api.Group("/payments")
		{
			payments.GET("", paymentController.GetPayments)
		}

		// Promotion routes
		promotions := api.Group("/promotions")
		{
			promotions.GET("", middlewares.AuthMiddleware(), promotionController.GetAllPromotions)
			promotions.POST("", middlewares.AuthMiddleware(), promotionController.CreatePromotion)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.GET("", middlewares.AuthMiddleware(), orderController.GetOrdersByUserId)
			orders.GET(":id", middlewares.AuthMiddleware(), orderController.GetOrderById)
			orders.POST("", middlewares.AuthMiddleware(), orderController.CreateOrder)
			orders.DELETE(":id", middlewares.AuthMiddleware(), orderController.DeleteOrder)
		}
	}

	return r
}
