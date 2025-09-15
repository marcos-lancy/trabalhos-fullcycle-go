package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"rate-limiter/internal/limiter"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	config := &limiter.Config{
		IPRequestsPerSecond:       1,
		IPBlockDurationMinutes:    1,
		TokenRequestsPerSecond:    3,
		TokenBlockDurationMinutes: 1,
	}

	rateLimiter := limiter.NewRateLimiter(storage, config)

	// Create test router
	router := gin.New()
	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		ip := c.ClientIP()
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

		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Test IP rate limiting
	t.Run("IP Rate Limiting", func(t *testing.T) {
		// Make first request (should be allowed)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		req.Header.Set("X-Real-IP", "192.168.1.1")

		w := &mockResponseWriter{}
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.statusCode, "First request should be allowed")

		// Make second request (should be blocked)
		req, _ = http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		req.Header.Set("X-Real-IP", "192.168.1.1")

		w = &mockResponseWriter{}
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.statusCode, "Second request should be blocked")
		assert.Contains(t, w.body.String(), "you have reached the maximum number of requests")
	})

	// Test Token rate limiting
	t.Run("Token Rate Limiting", func(t *testing.T) {
		// Make requests within token limit
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.2:12345"
			req.Header.Set("X-Real-IP", "192.168.1.2")
			req.Header.Set("API_KEY", "test-token")

			w := &mockResponseWriter{}
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.statusCode, "Request %d should be allowed", i+1)
		}

		// Make request that exceeds token limit
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		req.Header.Set("X-Real-IP", "192.168.1.2")
		req.Header.Set("API_KEY", "test-token")

		w := &mockResponseWriter{}
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.statusCode, "Request should be blocked")
		assert.Contains(t, w.body.String(), "you have reached the maximum number of requests")
	})

	// Test Token overrides IP
	t.Run("Token Overrides IP", func(t *testing.T) {
		// Reset storage
		storage.Delete(context.Background(), "ip:192.168.1.3")
		storage.Delete(context.Background(), "token:override-token")

		// Make requests with token (should use token limit, not IP limit)
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.3:12345"
			req.Header.Set("X-Real-IP", "192.168.1.3")
			req.Header.Set("API_KEY", "override-token")

			w := &mockResponseWriter{}
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.statusCode, "Request %d should be allowed with token", i+1)
		}

		// Now make request without token (should use IP limit - separate counter)
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		req.Header.Set("X-Real-IP", "192.168.1.3")

		w := &mockResponseWriter{}
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.statusCode, "First IP request should be allowed")

		// Second IP request should be blocked (IP limit is 2 req/s)
		req, _ = http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		req.Header.Set("X-Real-IP", "192.168.1.3")

		w = &mockResponseWriter{}
		router.ServeHTTP(w, req)

		// Debug output
		t.Logf("Response status: %d", w.statusCode)
		t.Logf("Response body: %s", w.body.String())

		assert.Equal(t, 429, w.statusCode, "Second IP request should be blocked")
	})
}

// mockResponseWriter implements http.ResponseWriter for testing
type mockResponseWriter struct {
	statusCode int
	body       bytes.Buffer
	headers    http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	if m.headers == nil {
		m.headers = make(http.Header)
	}
	return m.headers
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	return m.body.Write(data)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestRateLimiterLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	// Setup
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	config := &limiter.Config{
		IPRequestsPerSecond:       10,
		IPBlockDurationMinutes:    1,
		TokenRequestsPerSecond:    20,
		TokenBlockDurationMinutes: 1,
	}

	rateLimiter := limiter.NewRateLimiter(storage, config)
	ctx := context.Background()

	// Test concurrent requests
	concurrency := 50
	requestsPerGoroutine := 5

	results := make(chan bool, concurrency*requestsPerGoroutine)

	for i := 0; i < concurrency; i++ {
		go func(goroutineID int) {
			for j := 0; j < requestsPerGoroutine; j++ {
				ip := fmt.Sprintf("192.168.1.%d", goroutineID%10)
				token := fmt.Sprintf("token-%d", goroutineID%5)

				allowed, err := rateLimiter.CheckRequest(ctx, ip, token)
				results <- (allowed && err == nil)
			}
		}(i)
	}

	// Collect results
	successCount := 0
	totalRequests := concurrency * requestsPerGoroutine

	for i := 0; i < totalRequests; i++ {
		if <-results {
			successCount++
		}
	}

	// Verify that some requests were allowed and some were blocked
	// This is a probabilistic test - we expect some success and some failures
	assert.Greater(t, successCount, 0, "Some requests should be allowed")
	assert.Less(t, successCount, totalRequests, "Some requests should be blocked")

	t.Logf("Load test results: %d/%d requests allowed", successCount, totalRequests)
}
