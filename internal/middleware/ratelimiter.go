package middleware

import (
	"net/http"
	"sync"

	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware creates a rate limiter middleware that limits requests per IP address
func RateLimiterMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	// Map to store rate limiters for each IP address
	limiters := make(map[string]*rate.Limiter)
	var mu sync.Mutex
	fmt.Println("RateLimiterMiddleware")
	// Function to get the rate limiter for a given IP address
	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}

		return limiter
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
