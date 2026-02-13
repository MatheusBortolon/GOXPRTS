package storage

import (
	"context"
	"time"
)

// Storage defines the interface for rate limiter storage
type Storage interface {
	// Increment increments the counter for a given key
	// Returns the new count and an error if any
	Increment(ctx context.Context, key string) (int64, error)

	// Get retrieves the current count for a given key
	Get(ctx context.Context, key string) (int64, error)

	// SetExpiration sets the expiration time for a key
	SetExpiration(ctx context.Context, key string, expiration time.Duration) error

	// IsBlocked checks if a key is currently blocked
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Block blocks a key for a given duration
	Block(ctx context.Context, key string, duration time.Duration) error

	// Reset resets the counter for a key
	Reset(ctx context.Context, key string) error

	// Close closes the storage connection
	Close() error
}
