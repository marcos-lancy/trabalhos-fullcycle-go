package limiter

import (
	"context"
	"testing"
	"time"

	"rate-limiter/internal/storage"
)

func TestRateLimiter_CheckRequest_IP(t *testing.T) {
	storage := storage.NewMemoryStorage()
	config := &Config{
		IPRequestsPerSecond:    2,
		IPBlockDurationMinutes: 1,
	}
	
	limiter := NewRateLimiter(storage, config)
	ctx := context.Background()
	
	// Test normal requests within limit
	for i := 0; i < 2; i++ {
		allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !allowed {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}
	
	// Test request that exceeds limit
	allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err == nil {
		t.Error("Expected error for exceeding rate limit")
	}
	if allowed {
		t.Error("Expected request to be blocked")
	}
}

func TestRateLimiter_CheckRequest_Token(t *testing.T) {
	storage := storage.NewMemoryStorage()
	config := &Config{
		TokenRequestsPerSecond:    3,
		TokenBlockDurationMinutes: 1,
	}
	
	limiter := NewRateLimiter(storage, config)
	ctx := context.Background()
	
	// Test normal requests within limit
	for i := 0; i < 3; i++ {
		allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "token123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !allowed {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}
	
	// Test request that exceeds limit
	allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "token123")
	if err == nil {
		t.Error("Expected error for exceeding rate limit")
	}
	if allowed {
		t.Error("Expected request to be blocked")
	}
}

func TestRateLimiter_TokenOverridesIP(t *testing.T) {
	storage := storage.NewMemoryStorage()
	config := &Config{
		IPRequestsPerSecond:       1,
		IPBlockDurationMinutes:    1,
		TokenRequestsPerSecond:    5,
		TokenBlockDurationMinutes: 1,
	}
	
	limiter := NewRateLimiter(storage, config)
	ctx := context.Background()
	
	// Test that token allows more requests than IP limit
	for i := 0; i < 5; i++ {
		allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "token123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !allowed {
			t.Errorf("Expected request %d to be allowed with token", i+1)
		}
	}
	
	// Test that without token, IP limit applies
	allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error for first IP request, got %v", err)
	}
	if !allowed {
		t.Error("Expected first IP request to be allowed")
	}
	
	// Second IP request should be blocked
	allowed, err = limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err == nil {
		t.Error("Expected error for exceeding IP rate limit")
	}
	if allowed {
		t.Error("Expected second IP request to be blocked")
	}
}

func TestRateLimiter_GetRemainingRequests(t *testing.T) {
	storage := storage.NewMemoryStorage()
	config := &Config{
		IPRequestsPerSecond:    3,
		IPBlockDurationMinutes: 1,
	}
	
	limiter := NewRateLimiter(storage, config)
	ctx := context.Background()
	
	// Initially should have all requests available
	remaining, err := limiter.GetRemainingRequests(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if remaining != 3 {
		t.Errorf("Expected 3 remaining requests, got %d", remaining)
	}
	
	// Make one request
	_, err = limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Should have 2 remaining
	remaining, err = limiter.GetRemainingRequests(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if remaining != 2 {
		t.Errorf("Expected 2 remaining requests, got %d", remaining)
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	storage := storage.NewMemoryStorage()
	config := &Config{
		IPRequestsPerSecond:    1,
		IPBlockDurationMinutes: 1,
	}
	
	limiter := NewRateLimiter(storage, config)
	ctx := context.Background()
	
	// Make request to use up limit
	_, err := limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Second request should be blocked
	_, err = limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err == nil {
		t.Error("Expected error for exceeding rate limit")
	}
	
	// Reset the limit
	err = limiter.Reset(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Should be able to make request again
	allowed, err := limiter.CheckRequest(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !allowed {
		t.Error("Expected request to be allowed after reset")
	}
}

func TestConfig_LoadConfig(t *testing.T) {
	config := LoadConfig()
	
	// Test default values
	if config.IPRequestsPerSecond == 0 {
		t.Error("Expected IPRequestsPerSecond to have a default value")
	}
	if config.IPBlockDurationMinutes == 0 {
		t.Error("Expected IPBlockDurationMinutes to have a default value")
	}
	if config.TokenRequestsPerSecond == 0 {
		t.Error("Expected TokenRequestsPerSecond to have a default value")
	}
	if config.TokenBlockDurationMinutes == 0 {
		t.Error("Expected TokenBlockDurationMinutes to have a default value")
	}
}

func TestConfig_GetDurations(t *testing.T) {
	config := &Config{
		IPBlockDurationMinutes:    5,
		TokenBlockDurationMinutes: 10,
	}
	
	ipDuration := config.GetIPBlockDuration()
	if ipDuration != 5*time.Minute {
		t.Errorf("Expected IP block duration to be 5 minutes, got %v", ipDuration)
	}
	
	tokenDuration := config.GetTokenBlockDuration()
	if tokenDuration != 10*time.Minute {
		t.Errorf("Expected token block duration to be 10 minutes, got %v", tokenDuration)
	}
}
