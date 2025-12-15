package middleware

import (
	"go-backend-project/internal/ratelimit"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RateLimit(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		resultRateLimit, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "rate limiter error",
			})
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(resultRateLimit.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(resultRateLimit.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resultRateLimit.ResetAt, 10))

		if !resultRateLimit.Allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
			return
		}

		c.Next()
	}
}
