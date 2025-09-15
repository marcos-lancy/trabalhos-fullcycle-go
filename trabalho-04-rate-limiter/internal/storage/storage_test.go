package storage

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStorage_GetSet(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()
	
	// Test getting non-existent key
	data, err := storage.Get(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data != nil {
		t.Error("Expected nil for non-existent key")
	}
	
	// Test setting and getting data
	testData := &LimiterData{
		Count:     5,
		ResetTime: time.Now().Add(time.Minute),
		Blocked:   false,
		BlockUntil: time.Time{},
	}
	
	err = storage.Set(ctx, "test", testData, time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	retrieved, err := storage.Get(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected data to be retrieved")
	}
	if retrieved.Count != testData.Count {
		t.Errorf("Expected count %d, got %d", testData.Count, retrieved.Count)
	}
}

func TestMemoryStorage_Increment(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()
	
	// Test incrementing non-existent key
	count, err := storage.Increment(ctx, "test", time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
	
	// Test incrementing existing key
	count, err = storage.Increment(ctx, "test", time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestMemoryStorage_BlockUnblock(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()
	
	// Test setting block
	blockUntil := time.Now().Add(time.Minute)
	err := storage.SetBlock(ctx, "test", blockUntil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Test checking block
	blocked, blockTime, err := storage.IsBlocked(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !blocked {
		t.Error("Expected key to be blocked")
	}
	if blockTime.IsZero() {
		t.Error("Expected block time to be set")
	}
	
	// Test checking non-blocked key
	blocked, _, err = storage.IsBlocked(ctx, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if blocked {
		t.Error("Expected non-existent key to not be blocked")
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()
	
	// Set some data
	testData := &LimiterData{Count: 5}
	err := storage.Set(ctx, "test", testData, time.Minute)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Set a block
	err = storage.SetBlock(ctx, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Delete the key
	err = storage.Delete(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify data is gone
	data, err := storage.Get(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data != nil {
		t.Error("Expected data to be deleted")
	}
	
	// Verify block is gone
	blocked, _, err := storage.IsBlocked(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if blocked {
		t.Error("Expected block to be deleted")
	}
}

func TestMemoryStorage_Expiration(t *testing.T) {
	storage := NewMemoryStorage()
	ctx := context.Background()
	
	// Set data with short expiration
	testData := &LimiterData{Count: 5}
	err := storage.Set(ctx, "test", testData, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify data exists
	data, err := storage.Get(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data == nil {
		t.Error("Expected data to exist")
	}
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Verify data is gone
	data, err = storage.Get(ctx, "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if data != nil {
		t.Error("Expected data to be expired")
	}
}
