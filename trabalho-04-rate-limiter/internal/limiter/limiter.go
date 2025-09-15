package limiter

import (
	"context"
	"fmt"
	"time"

	"rate-limiter/internal/storage"
)

type RateLimiter struct {
	storage storage.Storage
	config  *Config
}

func NewRateLimiter(storage storage.Storage, config *Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) CheckRequest(ctx context.Context, ip, token string) (bool, error) {
	if blocked, blockUntil, err := rl.storage.IsBlocked(ctx, fmt.Sprintf("ip:%s", ip)); err != nil {
		return false, fmt.Errorf("failed to check IP block status: %w", err)
	} else if blocked {
		return false, fmt.Errorf("IP blocked until %v", blockUntil)
	}

	if token != "" {
		if blocked, blockUntil, err := rl.storage.IsBlocked(ctx, fmt.Sprintf("token:%s", token)); err != nil {
			return false, fmt.Errorf("failed to check token block status: %w", err)
		} else if blocked {
			return false, fmt.Errorf("token blocked until %v", blockUntil)
		}
	}

	var requestsPerSecond int
	var blockDuration time.Duration
	var keyPrefix string

	if token != "" {
		requestsPerSecond = rl.config.TokenRequestsPerSecond
		blockDuration = rl.config.GetTokenBlockDuration()
		keyPrefix = fmt.Sprintf("token:%s", token)
	} else {
		requestsPerSecond = rl.config.IPRequestsPerSecond
		blockDuration = rl.config.GetIPBlockDuration()
		keyPrefix = fmt.Sprintf("ip:%s", ip)
	}

	currentCount, err := rl.storage.Increment(ctx, keyPrefix, time.Second)
	if err != nil {
		return false, fmt.Errorf("failed to increment request count: %w", err)
	}

	if currentCount > int64(requestsPerSecond) {
		blockUntil := time.Now().Add(blockDuration)
		if err := rl.storage.SetBlock(ctx, keyPrefix, blockUntil); err != nil {
			return false, fmt.Errorf("failed to set block: %w", err)
		}
		return false, fmt.Errorf("rate limit exceeded, blocked until %v", blockUntil)
	}

	return true, nil
}

func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, ip, token string) (int64, error) {
	var requestsPerSecond int
	var keyPrefix string

	if token != "" {
		requestsPerSecond = rl.config.TokenRequestsPerSecond
		keyPrefix = fmt.Sprintf("token:%s", token)
	} else {
		requestsPerSecond = rl.config.IPRequestsPerSecond
		keyPrefix = fmt.Sprintf("ip:%s", ip)
	}

	data, err := rl.storage.Get(ctx, keyPrefix)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return int64(requestsPerSecond), nil
	}

	remaining := int64(requestsPerSecond) - data.Count
	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

func (rl *RateLimiter) Reset(ctx context.Context, ip, token string) error {
	var keyPrefix string

	if token != "" {
		keyPrefix = fmt.Sprintf("token:%s", token)
	} else {
		keyPrefix = fmt.Sprintf("ip:%s", ip)
	}

	return rl.storage.Delete(ctx, keyPrefix)
}
