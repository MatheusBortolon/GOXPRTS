package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis   RedisConfig
	Limiter LimiterConfig
	Server  ServerConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type LimiterConfig struct {
	IPRateLimit     int
	IPBlockTime     int
	TokenRateLimits map[string]TokenLimit
}

type TokenLimit struct {
	RPS       int
	BlockTime int
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Limiter: LimiterConfig{
			IPRateLimit:     getEnvAsInt("RATE_LIMIT_IP_RPS", 5),
			IPBlockTime:     getEnvAsInt("RATE_LIMIT_IP_BLOCK_TIME", 300),
			TokenRateLimits: parseTokenLimits(getEnv("RATE_LIMIT_TOKENS", "")),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func parseTokenLimits(tokens string) map[string]TokenLimit {
	limits := make(map[string]TokenLimit)
	if tokens == "" {
		return limits
	}

	tokenList := strings.Split(tokens, ",")
	for _, token := range tokenList {
		parts := strings.Split(strings.TrimSpace(token), ":")
		if len(parts) != 3 {
			continue
		}

		rps, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		blockTime, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		limits[parts[0]] = TokenLimit{
			RPS:       rps,
			BlockTime: blockTime,
		}
	}

	return limits
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
