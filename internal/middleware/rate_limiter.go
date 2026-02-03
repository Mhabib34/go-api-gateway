package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config
const (
	MaxRequest = 100
	WindowTime = time.Minute
)

// Middleware
func RateLimiterByIP(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := getClientIP(r)
			key := fmt.Sprintf("rate_limit:%s", ip)

			// Check rate limit
			count, err := rdb.Incr(r.Context(), key).Result() 
			if err != nil {
				http.Error(w, "Redis error", http.StatusInternalServerError)
				return
			}

			// Reset rate limit
			if count == 1 {
				rdb.Expire(r.Context(), key, WindowTime)
			}

			// Check rate limit
			if count > MaxRequest {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Serve request
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Get client IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	// Return client IP
	return ip
}
