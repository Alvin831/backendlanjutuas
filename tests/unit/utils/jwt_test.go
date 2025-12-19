package utils_test

import (
	"os"
	"testing"
	"time"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-min-32-chars")
	code := m.Run()
	os.Exit(code)
}

func TestGenerateToken(t *testing.T) {
	userID := "test-user-id"
	roleID := "test-role-id"
	permissions := []string{"read", "write", "admin"}

	tests := []struct {
		name        string
		userID      string
		roleID      string
		permissions []string
		wantErr     bool
	}{
		{
			name:        "Valid token generation",
			userID:      userID,
			roleID:      roleID,
			permissions: permissions,
			wantErr:     false,
		},
		{
			name:        "Empty user ID",
			userID:      "",
			roleID:      roleID,
			permissions: permissions,
			wantErr:     false, // JWT allows empty claims
		},
		{
			name:        "Nil permissions",
			userID:      userID,
			roleID:      roleID,
			permissions: nil,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tt.userID, tt.roleID, tt.permissions)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Contains(t, token, ".") // JWT has dots
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	userID := "test-user-id"
	roleID := "test-role-id"
	permissions := []string{"read", "write"}

	// Generate a valid token first
	validToken, err := utils.GenerateToken(userID, roleID, permissions)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "not.a.jwt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ParseToken(tt.token)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.UserID)
				assert.Equal(t, roleID, claims.Role)
				assert.Equal(t, permissions, claims.Permissions)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	userID := "test-user-id"
	roleID := "test-role-id"
	permissions := []string{"read"}

	// Generate token
	token, err := utils.GenerateToken(userID, roleID, permissions)
	assert.NoError(t, err)

	// Parse immediately (should be valid)
	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check expiration time is in the future
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestTokenRoundTrip(t *testing.T) {
	testCases := []struct {
		userID      string
		roleID      string
		permissions []string
	}{
		{"user1", "role1", []string{"perm1", "perm2"}},
		{"user2", "role2", []string{}},
		{"user3", "role3", nil},
	}

	for _, tc := range testCases {
		t.Run("RoundTrip_"+tc.userID, func(t *testing.T) {
			// Generate token
			token, err := utils.GenerateToken(tc.userID, tc.roleID, tc.permissions)
			assert.NoError(t, err)

			// Parse token
			claims, err := utils.ParseToken(token)
			assert.NoError(t, err)

			// Verify all fields match
			assert.Equal(t, tc.userID, claims.UserID)
			assert.Equal(t, tc.roleID, claims.Role)
			assert.Equal(t, tc.permissions, claims.Permissions)
		})
	}
}