package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/goxprts/ratelimiter/internal/storage"
)

type RateLimiter struct {
	storage         storage.Storage
	ipRateLimit     int
	ipBlockTime     time.Duration
	tokenRateLimits map[string]TokenConfig
}

type TokenConfig struct {
	RPS       int
	BlockTime time.Duration
}

type LimitResult struct {
	Allowed   bool
	Remaining int
	ResetTime time.Time
	Message   string
}

func NewRateLimiter(
	store storage.Storage,
	ipRateLimit int,
	ipBlockTime int,
	tokenLimits map[string]TokenConfig,
) *RateLimiter {
	return &RateLimiter{
		storage:         store,
		ipRateLimit:     ipRateLimit,
		ipBlockTime:     time.Duration(ipBlockTime) * time.Second,
		tokenRateLimits: tokenLimits,
	}
}

// Allow checks if a request is allowed based on IP or token
func (rl *RateLimiter) Allow(ctx context.Context, ip string, token string) (*LimitResult, error) {
	// Check if token is provided and has specific limits
	if token != "" {
		if tokenConfig, exists := rl.tokenRateLimits[token]; exists {
			return rl.checkLimit(ctx, fmt.Sprintf("token:%s", token), tokenConfig.RPS, tokenConfig.BlockTime)
		}
	}

	// Fall back to IP-based limiting
	return rl.checkLimit(ctx, fmt.Sprintf("ip:%s", ip), rl.ipRateLimit, rl.ipBlockTime)
}

func (rl *RateLimiter) checkLimit(ctx context.Context, key string, limit int, blockTime time.Duration) (*LimitResult, error) {
	// Check if the key is currently blocked
	blocked, err := rl.storage.IsBlocked(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to check if key is blocked: %w", err)
	}

	if blocked {
		return &LimitResult{
			Allowed:   false,
			Remaining: 0,
			ResetTime: time.Now().Add(blockTime),
			Message:   "you have reached the maximum number of requests or actions allowed within a certain time frame",
		}, nil
	}

	// Increment the counter
	count, err := rl.storage.Increment(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	// If this is the first request, set expiration to 1 second
	if count == 1 {
		if err := rl.storage.SetExpiration(ctx, key, time.Second); err != nil {
			return nil, fmt.Errorf("failed to set expiration: %w", err)
		}
	}

	// Check if limit is exceeded
	if count > int64(limit) {
		// Block the key
		if err := rl.storage.Block(ctx, key, blockTime); err != nil {
			return nil, fmt.Errorf("failed to block key: %w", err)
		}

		// Reset the counter
		if err := rl.storage.Reset(ctx, key); err != nil {
			return nil, fmt.Errorf("failed to reset counter: %w", err)
		}

		return &LimitResult{
			Allowed:   false,
			Remaining: 0,
			ResetTime: time.Now().Add(blockTime),
			Message:   "you have reached the maximum number of requests or actions allowed within a certain time frame",
		}, nil
	}

	remaining := limit - int(count)
	return &LimitResult{
		Allowed:   true,
		Remaining: remaining,
		ResetTime: time.Now().Add(time.Second),
		Message:   "",
	}, nil
}
