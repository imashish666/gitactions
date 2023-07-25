package middleware

import (
	"www-api/internal/logger"

	"github.com/gin-gonic/gin"
)

func LogRequest(log logger.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("request details", map[string]interface{}{
			"Address":  c.Request.RemoteAddr,
			"method":   c.Request.Method,
			"URL":      c.Request.URL,
			"Response": c.Writer.Status(),
		})
	}
}
