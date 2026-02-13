package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of storage.Storage
type MockStorage struct {
	mock.Mock
	counters map[string]int64
	blocked  map[string]bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counters: make(map[string]int64),
		blocked:  make(map[string]bool),
	}
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int64, error) {
	m.counters[key]++
	return m.counters[key], nil
}

func (m *MockStorage) Get(ctx context.Context, key string) (int64, error) {
	return m.counters[key], nil
}

func (m *MockStorage) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	return m.blocked[key], nil
}

func (m *MockStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	m.blocked[key] = true
	return nil
}

func (m *MockStorage) Reset(ctx context.Context, key string) error {
	m.counters[key] = 0
	return nil
}

func (m *MockStorage) Close() error {
	return nil
}

func TestRateLimiter_Allow_IPLimit(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := make(map[string]TokenConfig)

	limiter := NewRateLimiter(storage, 5, 300, tokenLimits)

	ctx := context.Background()
	ip := "192.168.1.1"

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		result, err := limiter.Allow(ctx, ip, "")
		assert.NoError(t, err)
		assert.True(t, result.Allowed)
		assert.Equal(t, 5-i-1, result.Remaining)
	}

	// 6th request should be blocked
	result, err := limiter.Allow(ctx, ip, "")
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame", result.Message)

	// 7th request should still be blocked
	result, err = limiter.Allow(ctx, ip, "")
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
}

func TestRateLimiter_Allow_TokenLimit(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := map[string]TokenConfig{
		"abc123": {
			RPS:       10,
			BlockTime: 300 * time.Second,
		},
	}

	limiter := NewRateLimiter(storage, 5, 300, tokenLimits)

	ctx := context.Background()
	ip := "192.168.1.1"
	token := "abc123"

	// First 10 requests should be allowed (token limit)
	for i := 0; i < 10; i++ {
		result, err := limiter.Allow(ctx, ip, token)
		assert.NoError(t, err)
		assert.True(t, result.Allowed)
		assert.Equal(t, 10-i-1, result.Remaining)
	}

	// 11th request should be blocked
	result, err := limiter.Allow(ctx, ip, token)
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
}

func TestRateLimiter_Allow_TokenOverridesIP(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := map[string]TokenConfig{
		"xyz789": {
			RPS:       100,
			BlockTime: 600 * time.Second,
		},
	}

	limiter := NewRateLimiter(storage, 5, 300, tokenLimits)

	ctx := context.Background()
	ip := "192.168.1.1"
	token := "xyz789"

	// Should use token limit (100) instead of IP limit (5)
	for i := 0; i < 10; i++ {
		result, err := limiter.Allow(ctx, ip, token)
		assert.NoError(t, err)
		assert.True(t, result.Allowed)
	}
}

func TestRateLimiter_Allow_DifferentIPs(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := make(map[string]TokenConfig)

	limiter := NewRateLimiter(storage, 5, 300, tokenLimits)

	ctx := context.Background()

	// IP1 uses 5 requests
	for i := 0; i < 5; i++ {
		result, err := limiter.Allow(ctx, "192.168.1.1", "")
		assert.NoError(t, err)
		assert.True(t, result.Allowed)
	}

	// IP2 should still have 5 requests available
	result, err := limiter.Allow(ctx, "192.168.1.2", "")
	assert.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Equal(t, 4, result.Remaining)
}
