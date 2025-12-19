package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"uas_backend/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-min-32-chars")
	os.Setenv("APP_ENV", "test")
	code := m.Run()
	os.Exit(code)
}

func setupTestServer() *fiber.App {
	app := fiber.New()
	
	// Add basic routes for testing
	app.Post("/test/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
		}
		
		// Mock successful login
		if req.Username == "testuser" && req.Password == "testpass" {
			return c.JSON(fiber.Map{
				"message": "Login successful",
				"token":   "mock-jwt-token",
			})
		}
		
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	})
	
	app.Get("/test/protected", func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth != "Bearer mock-jwt-token" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}
		
		return c.JSON(fiber.Map{"message": "Access granted"})
	})
	
	return app
}

func TestIntegration_LoginFlow(t *testing.T) {
	app := setupTestServer()
	
	// Test successful login
	loginReq := model.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}
	reqBody, _ := json.Marshal(loginReq)
	
	req := httptest.NewRequest("POST", "/test/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Parse response
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Equal(t, "Login successful", response["message"])
	assert.Equal(t, "mock-jwt-token", response["token"])
}

func TestIntegration_ProtectedRoute(t *testing.T) {
	app := setupTestServer()
	
	// Test without token
	req := httptest.NewRequest("GET", "/test/protected", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	
	// Test with valid token
	req = httptest.NewRequest("GET", "/test/protected", nil)
	req.Header.Set("Authorization", "Bearer mock-jwt-token")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Parse response
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Equal(t, "Access granted", response["message"])
}

func TestIntegration_InvalidLogin(t *testing.T) {
	app := setupTestServer()
	
	// Test with wrong credentials
	loginReq := model.LoginRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}
	reqBody, _ := json.Marshal(loginReq)
	
	req := httptest.NewRequest("POST", "/test/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	
	// Parse response
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Equal(t, "Invalid credentials", response["error"])
}