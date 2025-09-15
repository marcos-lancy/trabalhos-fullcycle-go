package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host, port, password string, db int) (*RedisStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: rdb}, nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (*LimiterData, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data LimiterData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal limiter data: %w", err)
	}

	return &data, nil
}

func (r *RedisStorage) Set(ctx context.Context, key string, data *LimiterData, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal limiter data: %w", err)
	}

	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.client.Pipeline()
	
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

func (r *RedisStorage) SetBlock(ctx context.Context, key string, blockUntil time.Time) error {
	blockKey := fmt.Sprintf("block:%s", key)
	duration := time.Until(blockUntil)
	
	if duration <= 0 {
		return nil
	}

	return r.client.Set(ctx, blockKey, "blocked", duration).Err()
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, time.Time, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	
	ttl, err := r.client.TTL(ctx, blockKey).Result()
	if err != nil {
		return false, time.Time{}, err
	}
	
	if ttl <= 0 {
		return false, time.Time{}, nil
	}

	blockUntil := time.Now().Add(ttl)
	return true, blockUntil, nil
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	blockKey := fmt.Sprintf("block:%s", key)
	
	pipe := r.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, blockKey)
	
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
