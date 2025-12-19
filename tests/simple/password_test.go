package simple_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test HashPassword function
func TestHashPassword(t *testing.T) {
	password := "123456"
	
	hash, err := utils.HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash) // Hash should be different from original
}

// Test CheckPasswordHash function
func TestCheckPasswordHash(t *testing.T) {
	password := "123456"
	hash, _ := utils.HashPassword(password)
	
	// Test correct password
	result := utils.CheckPasswordHash(password, hash)
	assert.True(t, result)
	
	// Test wrong password
	result = utils.CheckPasswordHash("wrongpassword", hash)
	assert.False(t, result)
}