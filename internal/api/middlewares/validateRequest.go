package middlewares

import (
	"net/url"
	"strings"

	"github.com/Ritika-Agrawal811/sheetdeck-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

/**
 * Middleware to validate the Origin and Referer headers against allowed origins
 * @param allowedOrigins []string - List of allowed origins
 * @return gin.HandlerFunc - Middleware function
 */
func ValidateRequestMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request's origin is allowed
		origin := c.GetHeader("Origin")
		if origin != "" && !isOriginAllowed(origin, allowedOrigins) {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Origin not allowed",
				"code":  "ORIGIN_NOT_ALLOWED",
			})
			return
		}

		// if origin is empty, check for referer
		referer := c.GetHeader("Referer")
		if origin == "" {
			if referer != "" && !isOriginAllowed(referer, allowedOrigins) {
				c.AbortWithStatusJSON(403, gin.H{
					"error": "Referer not allowed",
					"code":  "REFERER_NOT_ALLOWED",
				})
				return
			}
		}

		env := utils.GetEnv("ENV", "TEST")
		if env == "PROD" && origin == "" && referer == "" {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Origin or Referer header required in PROD environments",
				"code":  "MISSING_ORIGIN_REFERER",
			})
			return
		}

		c.Next()
	}
}

/**
 * Helper function to check if the origin is in the list of allowed origins
 * @param origin string - Origin from request
 * @param allowedOrigins []string - List of allowed origins
 * @return bool - true if allowed, false otherwise
 */
func isOriginAllowed(origin string, allowedOrigins []string) bool {

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	originSchemeHost := parsedOrigin.Scheme + "://" + parsedOrigin.Host

	for _, allowed := range allowedOrigins {
		parsedAllowed, err := url.Parse(allowed)
		if err != nil {
			continue
		}

		allowedSchemeHost := parsedAllowed.Scheme + "://" + parsedAllowed.Host

		if strings.EqualFold(originSchemeHost, allowedSchemeHost) {
			return true
		}
	}

	return false
}
