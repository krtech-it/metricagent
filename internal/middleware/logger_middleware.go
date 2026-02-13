package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		logger.Info("HTTP request",
			zap.String("URI", c.Request.RequestURI),
			zap.String("Method", c.Request.Method),
			zap.Duration("Latency", latency),
			zap.Int("StatusCode", c.Writer.Status()),
			zap.Int("Size", c.Writer.Size()),
		)
	}
}
