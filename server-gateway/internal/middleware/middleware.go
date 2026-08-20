package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger логирует входящие HTTP-запросы через zap.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		log.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// Recovery восстанавливает паники и отвечает 500.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)

				c.AbortWithStatusJSON(500, gin.H{
					"error": "internal server error",
					"code":  500,
				})
			}
		}()

		c.Next()
	}
}
