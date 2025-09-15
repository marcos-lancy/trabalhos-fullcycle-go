package main

import (
	"log"
	"os"
	"strconv"

	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Printf("Warning: Could not load config.env file: %v", err)
	}

	config := limiter.LoadConfig()

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := getEnvInt("REDIS_DB", 0)

	storage, err := storage.NewRedisStorage(redisHost, redisPort, redisPassword, redisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer storage.Close()

	rateLimiter := limiter.NewRateLimiter(storage, config)
	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"message": "Rate limiter is working",
		})
	})

	router.GET("/api/data", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is protected data",
			"timestamp": c.GetHeader("X-Request-Time"),
		})
	})

	router.POST("/api/data", func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message": "Data created successfully",
			"timestamp": c.GetHeader("X-Request-Time"),
		})
	})

	router.GET("/api/rate-limit/status", func(c *gin.Context) {
		ip := c.ClientIP()
		apiKey := c.GetHeader("API_KEY")
		
		remaining, err := rateLimiter.GetRemainingRequests(c.Request.Context(), ip, apiKey)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get rate limit status"})
			return
		}

		c.JSON(200, gin.H{
			"ip": ip,
			"has_token": apiKey != "",
			"remaining_requests": remaining,
		})
	})

	router.POST("/api/rate-limit/reset", func(c *gin.Context) {
		ip := c.ClientIP()
		apiKey := c.GetHeader("API_KEY")
		
		if err := rateLimiter.Reset(c.Request.Context(), ip, apiKey); err != nil {
			c.JSON(500, gin.H{"error": "Failed to reset rate limit"})
			return
		}

		c.JSON(200, gin.H{
			"message": "Rate limit reset successfully",
		})
	})

	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Printf("Rate limiter configuration:")
	log.Printf("  IP requests per second: %d", config.IPRequestsPerSecond)
	log.Printf("  IP block duration: %d minutes", config.IPBlockDurationMinutes)
	log.Printf("  Token requests per second: %d", config.TokenRequestsPerSecond)
	log.Printf("  Token block duration: %d minutes", config.TokenBlockDurationMinutes)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
