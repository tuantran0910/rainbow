package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/internal/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", controllers.NewBaseController().HomePage)
	r.GET("/health", controllers.NewBaseController().HealthCheck)

	return r
}
