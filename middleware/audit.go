package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	Response    interface{}            `json:"response,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
}

var auditLogger *log.Logger

func init() {
	// Buat folder logs jika belum ada
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		return
	}

	// Buat file audit log dengan nama berdasarkan tanggal
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join("logs", fmt.Sprintf("audit-%s.log", today))
	
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open audit log file: %v", err)
		return
	}

	auditLogger = log.New(file, "", 0)
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

// ====================================================================
// AUDIT MIDDLEWARE
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

		// Hitung durasi
		duration := time.Since(start)

		// Ambil user info dari context
		userID := "anonymous"
		role := "guest"
		
		if uid := c.Locals("user_id"); uid != nil {
			userID = uid.(string)
		}
		if r := c.Locals("role"); r != nil {
			role = r.(string)
		}

		// Buat audit log
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
			Status:    c.Response().StatusCode(),
			Duration:  duration.String(),
		}

		// Tambahkan request body jika ada
		if requestBody != nil {
			auditLog.RequestBody = requestBody
		}

		// Tambahkan headers penting
		auditLog.Headers = map[string]interface{}{
			"content-type":   c.Get("Content-Type"),
			"authorization":  maskAuthHeader(c.Get("Authorization")),
			"x-forwarded-for": c.Get("X-Forwarded-For"),
		}

		// Log hanya untuk endpoint penting atau error
		if shouldAuditLog(c.Path(), c.Response().StatusCode()) {
			writeAuditLog(auditLog)
		}

		return err
	}
}

// Helper functions
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
	// Extract resource dari path
	// /v1/users -> users
	// /auth/login -> auth
	// /roles/123 -> roles
	
	if path == "/" {
		return "root"
	}
	
	parts := []rune(path)
	if len(parts) > 0 && parts[0] == '/' {
		parts = parts[1:]
	}
	
	pathStr := string(parts)
	if pathStr == "" {
		return "root"
	}
	
	// Ambil bagian pertama setelah /
	for i, char := range pathStr {
		if char == '/' {
			return pathStr[:i]
		}
	}
	
	return pathStr
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
	}
	
	for _, importantPath := range importantPaths {
		if len(path) >= len(importantPath) && path[:len(importantPath)] == importantPath {
			return true
		}
	}
	
	return false
}