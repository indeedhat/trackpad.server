package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	ServerPort        = "SERVER_PORT"
	ConnetPass        = "CONNECT_PASS"
	DiscoveryInterval = "DISCOVERY_INTERVAL"
)

// Load envars from .env file
func Load() error {
	return godotenv.Load()
}

// Get environment variable by key, with optional fallback
func Get(key string, fallback ...string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return ""
}

// GetInt environment variable by key, with optional follback
func GetInt(key string, fallback ...int) int {
	val := os.Getenv(key)
	if val != "" {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return 0
}
