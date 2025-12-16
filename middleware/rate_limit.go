package middleware

import (
	"fmt"
	"sync"
	"time"
	"uas_backend/app/model"

	"github.com/gofiber/fiber/v2"
)

// ====================================================================
// RATE LIMITING — Mencegah abuse per user/permission
// ====================================================================

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

var rateLimiter = &RateLimiter{
	requests: make(map[string][]time.Time),
}

// Cleanup old requests setiap 5 menit
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			rateLimiter.cleanup()
		}
	}()
}

func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	cutoff := time.Now().Add(-time.Hour) // Hapus request lebih dari 1 jam
	
	for key, requests := range rl.requests {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(cutoff) {
				validRequests = append(validRequests, req)
			}
		}
		
		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}

func (rl *RateLimiter) isAllowed(key string, limit int, window time.Duration) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-window)
	
	// Ambil requests untuk key ini
	requests := rl.requests[key]
	
	// Filter requests yang masih dalam window
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	
	// Cek apakah masih dalam limit
	if len(validRequests) >= limit {
		return false
	}
	
	// Tambah request baru
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	
	return true
}

// ====================================================================
// RATE LIMIT MIDDLEWARE
// ====================================================================

// RateLimitByUser - Limit berdasarkan user ID
func RateLimitByUser(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			// Jika tidak ada user_id, skip rate limiting
			return c.Next()
		}
		
		key := fmt.Sprintf("user:%s", userID.(string))
		
		if !rateLimiter.isAllowed(key, limit, window) {
			return c.Status(429).JSON(model.WebResponse{
				Code:    429,
				Status:  "error",
				Message: fmt.Sprintf("Rate limit exceeded. Max %d requests per %v", limit, window),
			})
		}
		
		return c.Next()
	}
}

// RateLimitByPermission - Limit berdasarkan permission
func RateLimitByPermission(permission string, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Next()
		}
		
		key := fmt.Sprintf("user:%s:perm:%s", userID.(string), permission)
		
		if !rateLimiter.isAllowed(key, limit, window) {
			return c.Status(429).JSON(model.WebResponse{
				Code:    429,
				Status:  "error",
				Message: fmt.Sprintf("Rate limit exceeded for permission '%s'. Max %d requests per %v", permission, limit, window),
			})
		}
		
		return c.Next()
	}
}

// RateLimitByIP - Limit berdasarkan IP address
func RateLimitByIP(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("ip:%s", ip)
		
		if !rateLimiter.isAllowed(key, limit, window) {
			return c.Status(429).JSON(model.WebResponse{
				Code:    429,
				Status:  "error",
				Message: fmt.Sprintf("Rate limit exceeded from IP %s. Max %d requests per %v", ip, limit, window),
			})
		}
		
		return c.Next()
	}
}