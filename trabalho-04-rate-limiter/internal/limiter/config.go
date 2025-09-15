package limiter

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	IPRequestsPerSecond    int
	IPBlockDurationMinutes int
	TokenRequestsPerSecond    int
	TokenBlockDurationMinutes int
}

func LoadConfig() *Config {
	config := &Config{
		IPRequestsPerSecond:        getEnvInt("RATE_LIMIT_IP_REQUESTS_PER_SECOND", 5),
		IPBlockDurationMinutes:     getEnvInt("RATE_LIMIT_IP_BLOCK_DURATION_MINUTES", 5),
		TokenRequestsPerSecond:     getEnvInt("RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND", 10),
		TokenBlockDurationMinutes:  getEnvInt("RATE_LIMIT_TOKEN_BLOCK_DURATION_MINUTES", 5),
	}
	
	return config
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) GetIPBlockDuration() time.Duration {
	return time.Duration(c.IPBlockDurationMinutes) * time.Minute
}

func (c *Config) GetTokenBlockDuration() time.Duration {
	return time.Duration(c.TokenBlockDurationMinutes) * time.Minute
}
