package middleware

import (
	"san/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		// Calculate latency
		latency := time.Since(start)

		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("| %3d | %13v | %15s | %-7s %s | %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}
