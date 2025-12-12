package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf(`{
			"time": "%s",
			"method": "%s",
			"path": "%s",
			"status": %d,
			"latency_ms": %d,
			"client_ip": "%s",
			"user_agent": "%s"
		}`,
			start.Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start).Milliseconds(),
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}
