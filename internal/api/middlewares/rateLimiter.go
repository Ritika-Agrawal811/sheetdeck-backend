package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiterConfig struct {
	RequestsPerMinute int       // Number of requests made by this IP in current window
	ResetTime         time.Time // Time when the rate limit resets
}

type RateLimiterStore struct {
	clients map[string]*RateLimiterConfig // Map IP address -> rate limiter config
	mutex   sync.RWMutex
	limit   int
	window  time.Duration
}

func NewRateLimiterStore(limit int, window time.Duration) *RateLimiterStore {
	rl := &RateLimiterStore{
		clients: make(map[string]*RateLimiterConfig),
		limit:   limit,
		window:  window,
	}

	rl.startCleanupRoutine()

	return rl
}

/**
 * Middleware to enforce rate limiting based on client IP
 * @return gin.HandlerFunc - Gin middleware function
 */
func (rl *RateLimiterStore) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP address
		clientIP := c.ClientIP()

		if !rl.isAllowed(clientIP) {
			// Rate limit exceeded - return 429 Too Many Requests
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})

			return
		}

		c.Next()
	}
}

/**
 * Check if a request from the given IP is allowed under the rate limit
 * @param clientIP string - Client's IP address
 * @return bool - true if allowed, false if rate limit exceeded
 */
func (rl *RateLimiterStore) isAllowed(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	info, exists := rl.clients[clientIP]
	if !exists || now.After(info.ResetTime) {
		// New client or window has expired - reset count and window
		rl.clients[clientIP] = &RateLimiterConfig{
			RequestsPerMinute: 1,
			ResetTime:         now.Add(rl.window),
		}
		return true
	}

	if info.RequestsPerMinute < rl.limit {
		// Within limit - increment count
		info.RequestsPerMinute++
		return true
	}

	// Exceeded limit
	return false

}

/**
 * Cleanup old entries from the rate limiter store
 */
func (rl *RateLimiterStore) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for ip, info := range rl.clients {
		if now.After(info.ResetTime) {
			delete(rl.clients, ip)
		}
	}
}

/**
 * Start a goroutine to periodically clean up old entries
 */
func (rl *RateLimiterStore) startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			rl.cleanup()
		}
	}()
}
