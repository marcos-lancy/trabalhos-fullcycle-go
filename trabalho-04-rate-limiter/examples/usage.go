package main

import (
	"log"
	"net/http"
	"time"

	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
)

// Example of how to integrate the rate limiter into your own application
func main() {
	// Initialize Redis storage
	storage, err := storage.NewRedisStorage("localhost", "6379", "", 0)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer storage.Close()

	// Load configuration
	config := limiter.LoadConfig()

	// Initialize rate limiter
	rateLimiter := limiter.NewRateLimiter(storage, config)

	// Create Gin router
	router := gin.Default()

	// Add rate limiting middleware
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Your application routes
	router.GET("/api/users", getUsers)
	router.POST("/api/users", createUser)
	router.GET("/api/products", getProducts)
	router.POST("/api/orders", createOrder)

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Example handlers
func getUsers(c *gin.Context) {
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	c.JSON(http.StatusOK, gin.H{
		"users": []gin.H{
			{"id": 1, "name": "John Doe"},
			{"id": 2, "name": "Jane Smith"},
		},
	})
}

func createUser(c *gin.Context) {
	// Simulate user creation
	time.Sleep(200 * time.Millisecond)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": 123,
	})
}

func getProducts(c *gin.Context) {
	// Simulate product listing
	time.Sleep(150 * time.Millisecond)

	c.JSON(http.StatusOK, gin.H{
		"products": []gin.H{
			{"id": 1, "name": "Product 1", "price": 29.99},
			{"id": 2, "name": "Product 2", "price": 39.99},
		},
	})
}

func createOrder(c *gin.Context) {
	// Simulate order creation
	time.Sleep(300 * time.Millisecond)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order created successfully",
		"order_id": 456,
	})
}
