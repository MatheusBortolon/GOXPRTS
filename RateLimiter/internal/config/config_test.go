package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check default values
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, "6379", cfg.Redis.Port)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, 5, cfg.Limiter.IPRateLimit)
	assert.Equal(t, 300, cfg.Limiter.IPBlockTime)
	assert.Equal(t, "8080", cfg.Server.Port)
}

func TestLoad_CustomValues(t *testing.T) {
	// Set environment variables
	os.Setenv("REDIS_HOST", "redis-server")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("RATE_LIMIT_IP_RPS", "10")
	os.Setenv("RATE_LIMIT_IP_BLOCK_TIME", "600")
	os.Setenv("SERVER_PORT", "9090")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check custom values
	assert.Equal(t, "redis-server", cfg.Redis.Host)
	assert.Equal(t, "6380", cfg.Redis.Port)
	assert.Equal(t, 1, cfg.Redis.DB)
	assert.Equal(t, 10, cfg.Limiter.IPRateLimit)
	assert.Equal(t, 600, cfg.Limiter.IPBlockTime)
	assert.Equal(t, "9090", cfg.Server.Port)

	// Cleanup
	os.Clearenv()
}

func TestParseTokenLimits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "single token",
			input:    "abc123:10:300",
			expected: 1,
		},
		{
			name:     "multiple tokens",
			input:    "abc123:10:300,xyz789:100:600",
			expected: 2,
		},
		{
			name:     "invalid format",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTokenLimits(tt.input)
			assert.Equal(t, tt.expected, len(result))
		})
	}
}

func TestParseTokenLimits_Values(t *testing.T) {
	input := "abc123:10:300,xyz789:100:600"
	result := parseTokenLimits(input)

	assert.Equal(t, 2, len(result))

	token1, exists := result["abc123"]
	assert.True(t, exists)
	assert.Equal(t, 10, token1.RPS)
	assert.Equal(t, 300, token1.BlockTime)

	token2, exists := result["xyz789"]
	assert.True(t, exists)
	assert.Equal(t, 100, token2.RPS)
	assert.Equal(t, 600, token2.BlockTime)
}

func TestRedisConfig_Address(t *testing.T) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: "6379",
	}

	assert.Equal(t, "localhost:6379", cfg.Address())
}
