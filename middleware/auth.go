package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

// ====================================================================
// PERMISSION CACHE — In-memory cache untuk permissions
// ====================================================================
type PermissionCache struct {
	cache map[string]CacheEntry
	mutex sync.RWMutex
}

type CacheEntry struct {
	Permissions []string
	ExpiresAt   time.Time
}

var permCache = &PermissionCache{
	cache: make(map[string]CacheEntry),
}

// Cache permissions untuk 5 menit
func (pc *PermissionCache) Set(userID string, permissions []string) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	
	pc.cache[userID] = CacheEntry{
		Permissions: permissions,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}
}

func (pc *PermissionCache) Get(userID string) ([]string, bool) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()
	
	entry, exists := pc.cache[userID]
	if !exists || time.Now().After(entry.ExpiresAt) {
		// Hapus entry yang expired
		if exists {
			delete(pc.cache, userID)
		}
		return nil, false
	}
	
	return entry.Permissions, true
}

// ====================================================================
// RATE LIMITING — Mencegah abuse per user/IP
// ====================================================================
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

var rateLimiter = &RateLimiter{
	requests: make(map[string][]time.Time),
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

// ====================================================================
// AUDIT LOGGING — Detailed logging untuk compliance
// ====================================================================
type AuditLog struct {
	Timestamp   string                 `json:"timestamp"`
	UserID      string                 `json:"user_id"`
	Role        string                 `json:"role"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	IP          string                 `json:"ip"`
	UserAgent   string                 `json:"user_agent"`
	Status      int                    `json:"status"`
	Duration    string                 `json:"duration"`
	RequestBody interface{}            `json:"request_body,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
}

var auditLogger *log.Logger

// ====================================================================
// INITIALIZATION — Setup cache cleanup dan audit logger
// ====================================================================
func init() {
	// Setup audit logger
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
	} else {
		today := time.Now().Format("2006-01-02")
		logFile := filepath.Join("logs", fmt.Sprintf("audit-%s.log", today))
		
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open audit log file: %v", err)
		} else {
			auditLogger = log.New(file, "", 0)
		}
	}

	// Cleanup expired cache entries dan rate limiter setiap 10 menit
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			// Cleanup permission cache
			permCache.mutex.Lock()
			now := time.Now()
			for userID, entry := range permCache.cache {
				if now.After(entry.ExpiresAt) {
					delete(permCache.cache, userID)
				}
			}
			permCache.mutex.Unlock()
			
			// Cleanup rate limiter
			rateLimiter.cleanup()
		}
	}()
}

// ====================================================================
// LOGGING HELPERS
// ====================================================================
func logAccessAttempt(userID, role, permission, endpoint, method, status string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[RBAC] %s | User: %s | Role: %s | Permission: %s | %s %s | Status: %s", 
		timestamp, userID, role, permission, method, endpoint, status)
}

func writeAuditLog(auditLog AuditLog) {
	if auditLogger == nil {
		return
	}

	logData, err := json.Marshal(auditLog)
	if err != nil {
		log.Printf("Failed to marshal audit log: %v", err)
		return
	}

	auditLogger.Println(string(logData))
}

func determineAction(method, path string) string {
	switch method {
	case "GET":
		return "READ"
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return method
	}
}

func determineResource(path string) string {
	if path == "/" {
		return "root"
	}
	
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[0]
	}
	
	return "unknown"
}

func maskAuthHeader(auth string) string {
	if auth == "" {
		return ""
	}
	
	if len(auth) > 20 {
		return auth[:10] + "***" + auth[len(auth)-7:]
	}
	
	return "***"
}

func shouldAuditLog(path string, status int) bool {
	// Log semua error responses
	if status >= 400 {
		return true
	}
	
	// Log endpoint penting
	importantPaths := []string{
		"/auth/login",
		"/auth/register", 
		"/users",
		"/roles",
		"/achievements",
	}
	
	for _, importantPath := range importantPaths {
		if strings.Contains(path, importantPath) {
			return true
		}
	}
	
	return false
}

// ====================================================================
// 1️⃣ AUTH REQUIRED — Validasi TOKEN dengan Rate Limiting dan Audit
// ====================================================================
func AuthRequired(c *fiber.Ctx) error {
	start := time.Now()
	
	// Rate limiting by IP (100 requests per minute)
	ip := c.IP()
	ipKey := fmt.Sprintf("ip:%s", ip)
	if !rateLimiter.isAllowed(ipKey, 100, time.Minute) {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - Rate Limited")
		return c.Status(429).JSON(model.WebResponse{
			Code:    429,
			Status:  "error",
			Message: "Rate limit exceeded. Too many requests from your IP",
		})
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - No Auth Header")
		writeAuditLog(createAuditLog(c, start, "anonymous", "guest", 401, nil))
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Missing authorization header",
		})
	}

	var token string
	
	// Handle both formats: "Bearer token" and "token"
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = authHeader
	}
	
	// Validate token is not empty
	if strings.TrimSpace(token) == "" {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - Empty Token")
		writeAuditLog(createAuditLog(c, start, "anonymous", "guest", 401, nil))
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid token format",
		})
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - Invalid Token")
		writeAuditLog(createAuditLog(c, start, "anonymous", "guest", 401, nil))
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid or expired token",
		})
	}

	// Rate limiting by user (50 requests per minute)
	userKey := fmt.Sprintf("user:%s", claims.UserID)
	if !rateLimiter.isAllowed(userKey, 50, time.Minute) {
		logAccessAttempt(claims.UserID, claims.Role, "auth", c.Path(), c.Method(), "FAILED - User Rate Limited")
		writeAuditLog(createAuditLog(c, start, claims.UserID, claims.Role, 429, nil))
		return c.Status(429).JSON(model.WebResponse{
			Code:    429,
			Status:  "error",
			Message: "Rate limit exceeded. Too many requests from your account",
		})
	}

	// Cache permissions untuk user ini
	permCache.Set(claims.UserID, claims.Permissions)

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("permissions", claims.Permissions)

	logAccessAttempt(claims.UserID, claims.Role, "auth", c.Path(), c.Method(), "SUCCESS")
	
	// Audit log untuk successful auth
	writeAuditLog(createAuditLog(c, start, claims.UserID, claims.Role, 200, nil))
	
	return c.Next()
}

// ====================================================================
// 2️⃣ PERMISSION REQUIRED — Enhanced dengan logging, cache, dan audit
// ====================================================================
func PermissionRequired(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		userID := c.Locals("user_id")
		role := c.Locals("role")
		
		userIDStr := ""
		roleStr := ""
		
		if userID != nil {
			userIDStr = userID.(string)
		}
		if role != nil {
			roleStr = role.(string)
		}

		// 1. Coba ambil dari cache dulu
		var userPerms []string
		if cachedPerms, found := permCache.Get(userIDStr); found {
			userPerms = cachedPerms
		} else {
			// 2. Jika tidak ada di cache, ambil dari context
			userPermsInterface := c.Locals("permissions")
			if userPermsInterface == nil {
				logAccessAttempt(userIDStr, roleStr, requiredPerm, c.Path(), c.Method(), "FAILED - No Permissions")
				writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 403, nil))
				return c.Status(403).JSON(model.WebResponse{
					Code:    403,
					Status:  "error",
					Message: "No permissions found (Auth middleware might be missing)",
				})
			}

			// 3. Konversi interface{} ke []string dengan aman
			switch v := userPermsInterface.(type) {
			case []string:
				userPerms = v
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok {
						userPerms = append(userPerms, s)
					}
				}
			}
		}

		// 4. Cek apakah user punya izin yang dibutuhkan
		hasPermission := false
		for _, p := range userPerms {
			if p == requiredPerm {
				hasPermission = true
				break
			}
		}

		if hasPermission {
			logAccessAttempt(userIDStr, roleStr, requiredPerm, c.Path(), c.Method(), "SUCCESS")
			writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 200, nil))
			return c.Next()
		}

		// 5. Jika tidak ketemu, tolak akses
		logAccessAttempt(userIDStr, roleStr, requiredPerm, c.Path(), c.Method(), "FAILED - Access Denied")
		writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 403, nil))
		return c.Status(403).JSON(model.WebResponse{
			Code:    403,
			Status:  "error",
			Message: fmt.Sprintf("Access denied — missing permission: %s", requiredPerm),
		})
	}
}

// ====================================================================
// 3️⃣ MULTIPLE PERMISSIONS — User harus punya salah satu dari permissions
// ====================================================================
func AnyPermissionRequired(requiredPerms ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		userID := c.Locals("user_id")
		role := c.Locals("role")
		
		userIDStr := ""
		roleStr := ""
		
		if userID != nil {
			userIDStr = userID.(string)
		}
		if role != nil {
			roleStr = role.(string)
		}

		// Ambil permissions user
		var userPerms []string
		if cachedPerms, found := permCache.Get(userIDStr); found {
			userPerms = cachedPerms
		} else {
			userPermsInterface := c.Locals("permissions")
			if userPermsInterface == nil {
				logAccessAttempt(userIDStr, roleStr, strings.Join(requiredPerms, "|"), c.Path(), c.Method(), "FAILED - No Permissions")
				writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 403, nil))
				return c.Status(403).JSON(model.WebResponse{
					Code:    403,
					Status:  "error",
					Message: "No permissions found",
				})
			}

			switch v := userPermsInterface.(type) {
			case []string:
				userPerms = v
			case []interface{}:
				for _, item := range v {
					if s, ok := item.(string); ok {
						userPerms = append(userPerms, s)
					}
				}
			}
		}

		// Cek apakah user punya salah satu permission yang dibutuhkan
		for _, userPerm := range userPerms {
			for _, requiredPerm := range requiredPerms {
				if userPerm == requiredPerm {
					logAccessAttempt(userIDStr, roleStr, strings.Join(requiredPerms, "|"), c.Path(), c.Method(), "SUCCESS")
					writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 200, nil))
					return c.Next()
				}
			}
		}

		logAccessAttempt(userIDStr, roleStr, strings.Join(requiredPerms, "|"), c.Path(), c.Method(), "FAILED - Access Denied")
		writeAuditLog(createAuditLog(c, start, userIDStr, roleStr, 403, nil))
		return c.Status(403).JSON(model.WebResponse{
			Code:    403,
			Status:  "error",
			Message: fmt.Sprintf("Access denied — missing any of permissions: %s", strings.Join(requiredPerms, ", ")),
		})
	}
}

// ====================================================================
// 4️⃣ AUDIT MIDDLEWARE — Comprehensive logging
// ====================================================================
func AuditMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Simpan request body untuk logging (hati-hati dengan sensitive data)
		var requestBody interface{}
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			bodyBytes := c.Body()
			if len(bodyBytes) > 0 && len(bodyBytes) < 10000 { // Max 10KB untuk logging
				var body map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &body); err == nil {
					// Hapus field sensitive
					delete(body, "password")
					delete(body, "token")
					delete(body, "secret")
					requestBody = body
				}
			}
		}

		// Process request
		err := c.Next()

		// Ambil user info dari context
		userID := "anonymous"
		role := "guest"
		
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(string)
		}
		if r := c.Locals("role"); r != nil {
			role = r.(string)
		}

		// Log hanya untuk endpoint penting atau error
		if shouldAuditLog(c.Path(), c.Response().StatusCode()) {
			auditLog := createAuditLog(c, start, userID, role, c.Response().StatusCode(), requestBody)
			writeAuditLog(auditLog)
		}

		return err
	}
}

// ====================================================================
// HELPER FUNCTION — Create audit log
// ====================================================================
func createAuditLog(c *fiber.Ctx, start time.Time, userID, role string, status int, requestBody interface{}) AuditLog {
	duration := time.Since(start)

	auditLog := AuditLog{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000Z"),
		UserID:    userID,
		Role:      role,
		Action:    determineAction(c.Method(), c.Path()),
		Resource:  determineResource(c.Path()),
		Method:    c.Method(),
		Path:      c.Path(),
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent"),
		Status:    status,
		Duration:  duration.String(),
	}

	// Tambahkan request body jika ada
	if requestBody != nil {
		auditLog.RequestBody = requestBody
	}

	// Tambahkan headers penting
	auditLog.Headers = map[string]interface{}{
		"content-type":     c.Get("Content-Type"),
		"authorization":    maskAuthHeader(c.Get("Authorization")),
		"x-forwarded-for":  c.Get("X-Forwarded-For"),
	}

	return auditLog
}