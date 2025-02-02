package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/config"
	_ "github.com/tuantran0910/rainbow/docs"
	"github.com/tuantran0910/rainbow/internal/routes"
	"github.com/tuantran0910/rainbow/pkg/databases"
	"github.com/tuantran0910/rainbow/pkg/utils/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Rainbow API
//	@version		1.0
//	@description	This is the API documentation for the Rainbow API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	tuan.tran
//	@contact.email	tntuan0910@gmail.com

//	@license.name	MIT
//	@license.url	http://opensource.org/licenses/MIT

// @host		localhost:5000
// @BasePath	/api/v1
func main() {
	// Load the application configurations
	if err := config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("failed to load the application configurations: %v", err))
	}

	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to get the application configurations: %v", err))
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := logger.InitLogger(); err != nil {
		panic(fmt.Sprintf("failed to initialize the logger: %v", err))
	}

	// Initialize the database connection
	dbConnector := databases.NewPostgresConnector()
	db, err := dbConnector.GetInstance()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize the database connection: %v", err))
	}

	// Close the connection as the application stops
	defer dbConnector.Close()

	// Create a new router
	r := routes.NewRouter(db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	if err := r.Run(":" + cfg.ServerConfig.ServerPort); err != nil {
		panic(fmt.Sprintf("failed to run the server: %v", err))
	}
}
