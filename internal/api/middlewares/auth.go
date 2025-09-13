package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
 * Middleware to validate API key from request headers
 * @param validAPIKeys []string - List of valid API keys
 * @return gin.HandlerFunc - Gin middleware function
 */
func APIKeyMiddleware(validAPIKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// If no API key provided
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "MISSING_API_KEY",
			})

			return
		}

		// Validate API key
		if !isValidAPIKey(apiKey, validAPIKeys) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})

			return
		}

		// API key is valid, proceed to next handler
		c.Next()
	}
}

/**
 * Helper function to check if API key is valid
 * @param apiKey string - API key from request
 * @param validKeys []string - List of valid API keys
 * @return bool - true if valid, false otherwise
 */
func isValidAPIKey(apiKey string, validKeys []string) bool {
	for _, validKey := range validKeys {
		if strings.TrimSpace(apiKey) == strings.TrimSpace(validKey) {
			return true
		}
	}
	return false
}
