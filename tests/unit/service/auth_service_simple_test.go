package service_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test utility functions that don't require complex mocking

func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "Valid password match",
			password: "123456",
			expected: true,
		},
		{
			name:     "Invalid password",
			password: "wrongpassword",
			expected: false,
		},
	}

	// Generate hash for testing
	correctHash, err := utils.HashPassword("123456")
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckPasswordHash(tt.password, correctHash)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenGeneration(t *testing.T) {
	userID := "test-user"
	roleID := "test-role"
	permissions := []string{"read", "write"}

	// Test token generation
	token, err := utils.GenerateToken(userID, roleID, permissions)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test token parsing
	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, roleID, claims.Role)
	assert.Equal(t, permissions, claims.Permissions)
}

func TestResponseFormatting(t *testing.T) {
	// Test success response
	successResp := utils.SuccessResponse("Test success", 200, map[string]string{"key": "value"})
	assert.Equal(t, "Test success", successResp.Meta.Message)
	assert.Equal(t, 200, successResp.Meta.Code)
	assert.Equal(t, "success", successResp.Meta.Status)

	// Test error response
	errorResp := utils.ErrorResponse("Test error", 400, nil)
	assert.Equal(t, "Test error", errorResp.Meta.Message)
	assert.Equal(t, 400, errorResp.Meta.Code)
	assert.Equal(t, "error", errorResp.Meta.Status)
}