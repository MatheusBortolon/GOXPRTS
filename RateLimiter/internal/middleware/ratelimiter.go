package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/goxprts/ratelimiter/internal/limiter"
)

type RateLimiterMiddleware struct {
	limiter *limiter.RateLimiter
}

func NewRateLimiterMiddleware(limiter *limiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}

// Middleware returns an HTTP middleware function
func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address
		ip := extractIP(r)

		// Extract API key from header
		token := r.Header.Get("API_KEY")

		// Check rate limit
		ctx := context.Background()
		result, err := m.limiter.Allow(ctx, ip, token)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !result.Allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "` + result.Message + `"}`))
			return
		}

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Remaining", string(rune(result.Remaining)))

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

// extractIP extracts the IP address from the request
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
