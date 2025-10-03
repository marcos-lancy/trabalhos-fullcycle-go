package auction

import (
	"os"
	"testing"
	"time"
)

func TestGetAuctionDuration(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "Valid duration 5 minutes",
			envValue: "5m",
			expected: 5 * time.Minute,
		},
		{
			name:     "Valid duration 1 hour",
			envValue: "1h",
			expected: 1 * time.Hour,
		},
		{
			name:     "Valid duration 30 seconds",
			envValue: "30s",
			expected: 30 * time.Second,
		},
		{
			name:     "Invalid duration, should default to 5 minutes",
			envValue: "invalid",
			expected: 5 * time.Minute,
		},
		{
			name:     "Empty env var, should default to 5 minutes",
			envValue: "",
			expected: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("AUCTION_DURATION", tt.envValue)
			} else {
				os.Unsetenv("AUCTION_DURATION")
			}
			defer os.Unsetenv("AUCTION_DURATION")

			result := getAuctionDuration()
			if result != tt.expected {
				t.Errorf("getAuctionDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAuctionExpirationCalculation(t *testing.T) {
	os.Setenv("AUCTION_DURATION", "2m")
	defer os.Unsetenv("AUCTION_DURATION")

	now := time.Now()
	auctionDuration := getAuctionDuration()
	expirationTime := now.Add(-auctionDuration)

	// An auction created 3 minutes ago should be expired
	oldAuctionTime := now.Add(-3 * time.Minute)
	if !oldAuctionTime.Before(expirationTime) {
		t.Error("Auction from 3 minutes ago should be expired")
	}

	// An auction created 1 minute ago should not be expired
	recentAuctionTime := now.Add(-1 * time.Minute)
	if recentAuctionTime.Before(expirationTime) {
		t.Error("Auction from 1 minute ago should not be expired")
	}
}
