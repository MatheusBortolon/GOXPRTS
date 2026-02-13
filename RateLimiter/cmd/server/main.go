package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/goxprts/ratelimiter/internal/config"
	"github.com/goxprts/ratelimiter/internal/limiter"
	"github.com/goxprts/ratelimiter/internal/middleware"
	"github.com/goxprts/ratelimiter/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Redis storage
	redisStorage, err := storage.NewRedisStorage(
		cfg.Redis.Address(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	defer redisStorage.Close()

	// Convert token limits from config to limiter format
	tokenLimits := make(map[string]limiter.TokenConfig)
	for token, limit := range cfg.Limiter.TokenRateLimits {
		tokenLimits[token] = limiter.TokenConfig{
			RPS:       limit.RPS,
			BlockTime: time.Duration(limit.BlockTime) * time.Second,
		}
	}

	// Initialize rate limiter
	rateLimiter := limiter.NewRateLimiter(
		redisStorage,
		cfg.Limiter.IPRateLimit,
		cfg.Limiter.IPBlockTime,
		tokenLimits,
	)

	// Initialize middleware
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	// Create HTTP router
	mux := http.NewServeMux()

	// Add a simple test endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Request successful"}`))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	// Wrap the router with rate limiter middleware
	handler := rateLimiterMiddleware.Middleware(mux)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("IP Rate Limit: %d req/s, Block Time: %ds", cfg.Limiter.IPRateLimit, cfg.Limiter.IPBlockTime)
	log.Printf("Token Limits configured: %d tokens", len(cfg.Limiter.TokenRateLimits))

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
