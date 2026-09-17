package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// Lanjutkan proses ke handler / route berikutnya
		c.Next()
		// Hitung latency / waktu eksekusi
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		// Field log terstruktur
		entry := logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"ip":         clientIP,
			"method":     method,
			"path":       path,
		})
		// Tangkap error jika ada error dari Gin context
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
			return
		}
		// Tentukan level log berdasarkan HTTP Status Code
		switch {
		case statusCode >= 500:
			entry.Error("Server Internal Error")
		case statusCode >= 400:
			entry.Warn("Client Error")
		default:
			entry.Info("Request Handled")
		}
	}
}
