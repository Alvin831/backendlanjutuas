package middleware

import (
	"fmt"
	"log"
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

// Cleanup expired cache entries setiap 10 menit
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			permCache.mutex.Lock()
			now := time.Now()
			for userID, entry := range permCache.cache {
				if now.After(entry.ExpiresAt) {
					delete(permCache.cache, userID)
				}
			}
			permCache.mutex.Unlock()
		}
	}()
}

// ====================================================================
// LOGGING HELPER
// ====================================================================
func logAccessAttempt(userID, role, permission, endpoint, method, status string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[RBAC] %s | User: %s | Role: %s | Permission: %s | %s %s | Status: %s", 
		timestamp, userID, role, permission, method, endpoint, status)
}

// ====================================================================
// 1️⃣ AUTH REQUIRED — Validasi TOKEN
// ====================================================================
func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - No Auth Header")
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Missing authorization header",
		})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - Invalid Token Format")
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid token format",
		})
	}

	claims, err := utils.ParseToken(tokenParts[1])
	if err != nil {
		logAccessAttempt("unknown", "unknown", "auth", c.Path(), c.Method(), "FAILED - Invalid Token")
		return c.Status(401).JSON(model.WebResponse{
			Code:    401,
			Status:  "error",
			Message: "Invalid or expired token",
		})
	}

	// Cache permissions untuk user ini
	permCache.Set(claims.UserID, claims.Permissions)

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("permissions", claims.Permissions)

	logAccessAttempt(claims.UserID, claims.Role, "auth", c.Path(), c.Method(), "SUCCESS")
	return c.Next()
}

// ====================================================================
// 2️⃣ PERMISSION REQUIRED — Enhanced dengan logging dan cache
// ====================================================================
func PermissionRequired(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			return c.Next()
		}

		// 5. Jika tidak ketemu, tolak akses
		logAccessAttempt(userIDStr, roleStr, requiredPerm, c.Path(), c.Method(), "FAILED - Access Denied")
		return c.Status(403).JSON(model.WebResponse{
			Code:    403,
			Status:  "error",
			Message: fmt.Sprintf("Access denied — missing permission: %s", requiredPerm),
		})
	}
}

// ====================================================================
// 3️⃣ ROLE REQUIRED — Middleware untuk check role
// ====================================================================
func RoleRequired(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		if roleStr != requiredRole {
			logAccessAttempt(userIDStr, roleStr, fmt.Sprintf("role:%s", requiredRole), c.Path(), c.Method(), "FAILED - Wrong Role")
			return c.Status(403).JSON(model.WebResponse{
				Code:    403,
				Status:  "error",
				Message: fmt.Sprintf("Access denied — required role: %s, your role: %s", requiredRole, roleStr),
			})
		}

		logAccessAttempt(userIDStr, roleStr, fmt.Sprintf("role:%s", requiredRole), c.Path(), c.Method(), "SUCCESS")
		return c.Next()
	}
}

// ====================================================================
// 4️⃣ MULTIPLE PERMISSIONS — User harus punya salah satu dari permissions
// ====================================================================
func AnyPermissionRequired(requiredPerms ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
					return c.Next()
				}
			}
		}

		logAccessAttempt(userIDStr, roleStr, strings.Join(requiredPerms, "|"), c.Path(), c.Method(), "FAILED - Access Denied")
		return c.Status(403).JSON(model.WebResponse{
			Code:    403,
			Status:  "error",
			Message: fmt.Sprintf("Access denied — missing any of permissions: %s", strings.Join(requiredPerms, ", ")),
		})
	}
}
