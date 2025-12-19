package simple_test

import (
	"os"
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test GenerateToken function
func TestGenerateToken(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-min-32-chars")
	
	userID := "user123"
	role := "admin"
	permissions := []string{"read", "write"}

	token, err := utils.GenerateToken(userID, role, permissions)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".") // JWT contains dots
}

// Test ParseToken function with valid token
func TestParseTokenValid(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-min-32-chars")
	
	userID := "user123"
	role := "admin"
	permissions := []string{"read", "write"}

	// Generate token first
	token, err := utils.GenerateToken(userID, role, permissions)
	assert.NoError(t, err)

	// Parse the token
	claims, err := utils.ParseToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, permissions, claims.Permissions)
}

// Test ParseToken function with invalid token
func TestParseTokenInvalid(t *testing.T) {
	invalidToken := "invalid.token.here"

	claims, err := utils.ParseToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

// Test ParseToken function with empty token
func TestParseTokenEmpty(t *testing.T) {
	claims, err := utils.ParseToken("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}