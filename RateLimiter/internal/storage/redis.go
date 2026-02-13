package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{
		client: client,
	}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return count, nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (r *RedisStorage) SetExpiration(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	exists, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if key %s is blocked: %w", key, err)
	}
	return exists > 0, nil
}

func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	blockKey := fmt.Sprintf("block:%s", key)
	err := r.client.Set(ctx, blockKey, "1", duration).Err()
	if err != nil {
		return fmt.Errorf("failed to block key %s: %w", key, err)
	}
	return nil
}

func (r *RedisStorage) Reset(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to reset key %s: %w", key, err)
	}
	return nil
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
