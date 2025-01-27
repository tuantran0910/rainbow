package main

import (
	_ "github.com/tuantran0910/rainbow/docs"
	"github.com/tuantran0910/rainbow/internal/routes"
	"github.com/tuantran0910/rainbow/pkg/config"
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
	// Load the application configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading the application configurations: ", zap.Error(err))
	}

	// Initialize the database connection
	dbConnector := postgres.NewPostgresConnector(cfg)
	db, err := dbConnector.GetInstance()
	if err != nil {
		log.Fatal("Error initializing the database connection: ", zap.Error(err))
	}

	// Close the connection as the application stops
	defer dbConnector.Close()

	// Create a new router
	r := routes.NewRouter(db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	serverPort := cfg.ServerConfig.Port
	if serverPort == "" {
		serverPort = "5000"
	}
	r.Run(":" + serverPort)
}
