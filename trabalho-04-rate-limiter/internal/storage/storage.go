package storage

import (
	"context"
	"time"
)

type LimiterData struct {
	Count      int64     `json:"count"`
	ResetTime  time.Time `json:"reset_time"`
	Blocked    bool      `json:"blocked"`
	BlockUntil time.Time `json:"block_until"`
}

type Storage interface {
	Get(ctx context.Context, key string) (*LimiterData, error)
	Set(ctx context.Context, key string, data *LimiterData, expiration time.Duration) error
	Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
	SetBlock(ctx context.Context, key string, blockUntil time.Time) error
	IsBlocked(ctx context.Context, key string) (bool, time.Time, error)
	Delete(ctx context.Context, key string) error
	Close() error
}
