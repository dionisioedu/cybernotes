package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func gertLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if exists {
		return limiter
	}

	limiter = rate.NewLimiter(1, 3)
	visitors[ip] = limiter
	go func() {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		delete(visitors, ip)
		mu.Unlock()
	}()

	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := gertLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Muitas requisições. Tente novamente mais tarde"})
			c.Abort()
			return
		}

		c.Next()
	}
}
