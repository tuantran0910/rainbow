package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuantran0910/rainbow/pkg/utils/logger"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Get the logger
		log, err := logger.GetLogger()
		if err != nil {
			panic(err)
		}

		log.Sugar().Infof(
			"%s %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}
