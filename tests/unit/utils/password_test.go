package utils_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "123456",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Long password",
			password: "this_is_a_very_long_password_with_special_characters_!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := utils.HashPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash) // Hash should be different from original
				assert.True(t, len(hash) > 50) // bcrypt hashes are typically 60 chars
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "test123"
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "Wrong password",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid_hash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPasswordHashConsistency(t *testing.T) {
	password := "consistent_test"
	
	// Hash the same password multiple times
	hash1, err1 := utils.HashPassword(password)
	hash2, err2 := utils.HashPassword(password)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	
	// Hashes should be different (bcrypt uses salt)
	assert.NotEqual(t, hash1, hash2)
	
	// But both should validate correctly
	assert.True(t, utils.CheckPasswordHash(password, hash1))
	assert.True(t, utils.CheckPasswordHash(password, hash2))
}