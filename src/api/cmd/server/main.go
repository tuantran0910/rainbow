package main

import (
	"github.com/tuantran0910/rainbow/config"
	_ "github.com/tuantran0910/rainbow/docs"
	"github.com/tuantran0910/rainbow/internal/routes"
	"github.com/tuantran0910/rainbow/pkg/databases/postgres"
	"github.com/tuantran0910/rainbow/pkg/utils/logger"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var log = logger.GetLogger()

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
	// TODO: Set the Gin to release mode
	// gin.SetMode(gin.ReleaseMode)

	// Load the application configurations
	if err := config.LoadConfig(); err != nil {
		log.Error("Error loading the application configurations: ", zap.Error(err))
	}

	// Initialize the database connection
	dbConnector := postgres.NewPostgresConnector()
	db, err := dbConnector.GetInstance()
	if err != nil {
		log.Error("Error initializing the database connection: ", zap.Error(err))
	}

	// Close the connection as the application stops
	defer dbConnector.Close()

	// Create a new router
	r := routes.NewRouter(db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	cfg := config.GetConfig()
	if err := r.Run(":" + cfg.ServerConfig.ServerPort); err != nil {
		log.Error("Error running the server: ", zap.Error(err))
	}
}
