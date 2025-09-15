package storage

import (
	"context"
	"sync"
	"time"
)

type MemoryStorage struct {
	data   map[string]*LimiterData
	blocks map[string]time.Time
	mutex  sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:   make(map[string]*LimiterData),
		blocks: make(map[string]time.Time),
	}
}

func (m *MemoryStorage) Get(ctx context.Context, key string) (*LimiterData, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data, exists := m.data[key]
	if !exists {
		return nil, nil
	}

	return &LimiterData{
		Count:      data.Count,
		ResetTime:  data.ResetTime,
		Blocked:    data.Blocked,
		BlockUntil: data.BlockUntil,
	}, nil
}

func (m *MemoryStorage) Set(ctx context.Context, key string, data *LimiterData, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = &LimiterData{
		Count:      data.Count,
		ResetTime:  data.ResetTime,
		Blocked:    data.Blocked,
		BlockUntil: data.BlockUntil,
	}

	go func() {
		time.Sleep(expiration)
		m.mutex.Lock()
		delete(m.data, key)
		m.mutex.Unlock()
	}()

	return nil
}

func (m *MemoryStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data, exists := m.data[key]
	if !exists {
		data = &LimiterData{
			Count:     0,
			ResetTime: time.Now().Add(expiration),
		}
		m.data[key] = data
	}

	data.Count++

	go func() {
		time.Sleep(expiration)
		m.mutex.Lock()
		delete(m.data, key)
		m.mutex.Unlock()
	}()

	return data.Count, nil
}

func (m *MemoryStorage) SetBlock(ctx context.Context, key string, blockUntil time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.blocks[key] = blockUntil

	duration := time.Until(blockUntil)
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			m.mutex.Lock()
			delete(m.blocks, key)
			m.mutex.Unlock()
		}()
	}

	return nil
}

func (m *MemoryStorage) IsBlocked(ctx context.Context, key string) (bool, time.Time, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	blockUntil, exists := m.blocks[key]
	if !exists {
		return false, time.Time{}, nil
	}

	if time.Now().After(blockUntil) {
		go func() {
			m.mutex.Lock()
			delete(m.blocks, key)
			m.mutex.Unlock()
		}()
		return false, time.Time{}, nil
	}

	return true, blockUntil, nil
}

func (m *MemoryStorage) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
	delete(m.blocks, key)

	return nil
}

func (m *MemoryStorage) Close() error {
	return nil
}
