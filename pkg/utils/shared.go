package utils

import (
	"os"
)

/**
 * Returns environment variables with fallback.
 */
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
