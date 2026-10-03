// Package httpx holds the small HTTP helpers both services share.
package httpx

import (
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Error writes {"error": msg} with the given status and stops the handler chain.
func Error(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

// CORS allows the listed front-end origins to call the service with a bearer token.
func CORS(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && slices.Contains(origins, origin) {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Vary", "Origin")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RateLimit allows each client IP `burst` requests, refilled at one per `every`.
func RateLimit(burst int, every time.Duration) gin.HandlerFunc {
	type bucket struct {
		tokens float64
		last   time.Time
	}
	var (
		mu      sync.Mutex
		buckets = map[string]*bucket{}
	)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: float64(burst), last: now}
			buckets[ip] = b
			if len(buckets) > 10000 {
				for k, v := range buckets {
					if now.Sub(v.last) > time.Hour {
						delete(buckets, k)
					}
				}
			}
		}
		b.tokens = min(float64(burst), b.tokens+now.Sub(b.last).Seconds()/every.Seconds())
		b.last = now
		allowed := b.tokens >= 1
		if allowed {
			b.tokens--
		}
		mu.Unlock()
		if !allowed {
			Error(c, http.StatusTooManyRequests, "請求太頻繁，請稍後再試")
			return
		}
		c.Next()
	}
}

// NewEngine returns a gin engine with CORS and a /healthz route. Client IPs come from
// X-Forwarded-For only when the request arrives from one of trustedProxies.
func NewEngine(origins, trustedProxies []string) *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		panic(err)
	}
	r.Use(CORS(origins))
	r.GET("/healthz", Healthz)
	return r
}

func Healthz(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) }
