package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors     = make(map[string]*rate.Limiter)
	mu           sync.Mutex
	cleanupOnce  sync.Once
)

// RateLimit implements rate limiting middleware
func RateLimit() gin.HandlerFunc {
	// Start cleanup goroutine only once
	cleanupOnce.Do(func() {
		go func() {
			for {
				time.Sleep(time.Minute)
				mu.Lock()
				// Reset the map periodically to prevent memory leak
				visitors = make(map[string]*rate.Limiter)
				mu.Unlock()
			}
		}()
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Allow 10 requests per second with burst of 20
		limiter = rate.NewLimiter(10, 20)
		visitors[ip] = limiter
	}

	return limiter
}
