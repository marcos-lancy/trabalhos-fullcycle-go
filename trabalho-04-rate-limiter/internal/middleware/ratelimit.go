package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"rate-limiter/internal/limiter"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(rateLimiter *limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		ip := getClientIP(c)
		apiKey := c.GetHeader("API_KEY")
		
		allowed, err := rateLimiter.CheckRequest(ctx, ip, apiKey)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			c.Abort()
			return
		}

		remaining, err := rateLimiter.GetRemainingRequests(ctx, ip, apiKey)
		if err == nil {
			c.Header("X-RateLimit-Remaining", string(rune(remaining)))
		}

		c.Next()
	}
}

func getClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	ip := c.ClientIP()
	if ip == "" {
		return "127.0.0.1"
	}

	return ip
}
