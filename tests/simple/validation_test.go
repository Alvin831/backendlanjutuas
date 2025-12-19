package simple_test

import (
	"testing"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test ValidateEmail function
func TestValidateEmail(t *testing.T) {
	// Valid emails
	assert.True(t, utils.ValidateEmail("test@example.com"))
	assert.True(t, utils.ValidateEmail("user.name@domain.co.id"))
	assert.True(t, utils.ValidateEmail("admin123@test.org"))
	
	// Invalid emails
	assert.False(t, utils.ValidateEmail(""))
	assert.False(t, utils.ValidateEmail("invalid-email"))
	assert.False(t, utils.ValidateEmail("@domain.com"))
	assert.False(t, utils.ValidateEmail("user@"))
	assert.False(t, utils.ValidateEmail("user@domain"))
}

// Test ValidateUsername function
func TestValidateUsername(t *testing.T) {
	// Valid usernames
	assert.True(t, utils.ValidateUsername("admin"))
	assert.True(t, utils.ValidateUsername("user123"))
	assert.True(t, utils.ValidateUsername("test_user"))
	
	// Invalid usernames
	assert.False(t, utils.ValidateUsername("ab")) // too short
	assert.False(t, utils.ValidateUsername("")) // empty
	assert.False(t, utils.ValidateUsername("user-name")) // contains dash
	assert.False(t, utils.ValidateUsername("user name")) // contains space
	assert.False(t, utils.ValidateUsername("verylongusernamethatexceedslimit")) // too long
}

// Test ValidatePassword function
func TestValidatePassword(t *testing.T) {
	// Valid passwords
	assert.True(t, utils.ValidatePassword("123456"))
	assert.True(t, utils.ValidatePassword("password"))
	assert.True(t, utils.ValidatePassword("verylongpassword"))
	
	// Invalid passwords
	assert.False(t, utils.ValidatePassword("")) // empty
	assert.False(t, utils.ValidatePassword("12345")) // too short
	assert.False(t, utils.ValidatePassword("abc")) // too short
}

// Test IsEmptyString function
func TestIsEmptyString(t *testing.T) {
	// Empty strings
	assert.True(t, utils.IsEmptyString(""))
	assert.True(t, utils.IsEmptyString("   ")) // only spaces
	assert.True(t, utils.IsEmptyString("\t\n")) // only whitespace
	
	// Non-empty strings
	assert.False(t, utils.IsEmptyString("hello"))
	assert.False(t, utils.IsEmptyString(" hello ")) // has content
	assert.False(t, utils.IsEmptyString("123"))
}