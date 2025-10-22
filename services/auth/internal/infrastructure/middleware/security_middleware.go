package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// XSS Protection
		c.Header("X-XSS-Protection", "1; mode=block")
		// No sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Frame options
		c.Header("X-Frame-Options", "DENY")
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// RateLimit implements rate limiting
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
