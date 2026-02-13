package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goxprts/ratelimiter/internal/limiter"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
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

func TestRateLimiterMiddleware_AllowedRequest(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := make(map[string]limiter.TokenConfig)
	rl := limiter.NewRateLimiter(storage, 5, 300, tokenLimits)
	middleware := NewRateLimiterMiddleware(rl)

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestRateLimiterMiddleware_BlockedRequest(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := make(map[string]limiter.TokenConfig)
	rl := limiter.NewRateLimiter(storage, 2, 300, tokenLimits)
	middleware := NewRateLimiterMiddleware(rl)

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Make 2 requests (at limit)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3rd request should be blocked
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "you have reached the maximum number of requests")
}

func TestRateLimiterMiddleware_WithToken(t *testing.T) {
	storage := NewMockStorage()
	tokenLimits := map[string]limiter.TokenConfig{
		"test-token": {
			RPS:       10,
			BlockTime: 300 * time.Second,
		},
	}
	rl := limiter.NewRateLimiter(storage, 2, 300, tokenLimits)
	middleware := NewRateLimiterMiddleware(rl)

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Make 5 requests with token (should all pass because token limit is 10)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestExtractIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.1")
	req.RemoteAddr = "192.168.1.1:1234"

	ip := extractIP(req)
	assert.Equal(t, "203.0.113.1", ip)
}

func TestExtractIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "203.0.113.1")
	req.RemoteAddr = "192.168.1.1:1234"

	ip := extractIP(req)
	assert.Equal(t, "203.0.113.1", ip)
}

func TestExtractIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"

	ip := extractIP(req)
	assert.Equal(t, "192.168.1.1", ip)
}
