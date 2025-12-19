package simple_test

import (
	"testing"
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test helper functions for services

// Test ValidateUserInput function (helper for user service)
func TestValidateUserInput(t *testing.T) {
	// Valid user input
	validUser := model.User{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   "role123",
	}
	
	assert.True(t, validateUserInput(validUser))
	
	// Invalid user input - empty username
	invalidUser1 := model.User{
		Username: "",
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   "role123",
	}
	
	assert.False(t, validateUserInput(invalidUser1))
	
	// Invalid user input - invalid email
	invalidUser2 := model.User{
		Username: "testuser",
		Email:    "invalid-email",
		FullName: "Test User",
		RoleID:   "role123",
	}
	
	assert.False(t, validateUserInput(invalidUser2))
}

// Test CalculatePagination function (helper for pagination)
func TestCalculatePagination(t *testing.T) {
	// Test normal pagination
	page, limit, offset := calculatePagination(2, 10)
	assert.Equal(t, 2, page)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 10, offset) // (2-1) * 10
	
	// Test first page
	page, limit, offset = calculatePagination(1, 5)
	assert.Equal(t, 1, page)
	assert.Equal(t, 5, limit)
	assert.Equal(t, 0, offset) // (1-1) * 5
	
	// Test invalid page (should default to 1)
	page, limit, offset = calculatePagination(0, 10)
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, offset)
}

// Test CalculateTotalPages function
func TestCalculateTotalPages(t *testing.T) {
	// Exact division
	totalPages := calculateTotalPages(100, 10)
	assert.Equal(t, 10, totalPages)
	
	// With remainder
	totalPages = calculateTotalPages(105, 10)
	assert.Equal(t, 11, totalPages)
	
	// Zero total
	totalPages = calculateTotalPages(0, 10)
	assert.Equal(t, 0, totalPages)
	
	// Single item
	totalPages = calculateTotalPages(1, 10)
	assert.Equal(t, 1, totalPages)
}

// Test FormatUserResponse function (helper to format user data)
func TestFormatUserResponse(t *testing.T) {
	user := &model.User{
		ID:           "user123",
		Username:     "testuser",
		Email:        "test@example.com",
		FullName:     "Test User",
		PasswordHash: "hashedpassword",
		RoleID:       "role123",
		IsActive:     true,
	}
	
	formatted := formatUserResponse(user)
	
	// Should remove password hash
	assert.Empty(t, formatted.PasswordHash)
	
	// Should keep other fields
	assert.Equal(t, "user123", formatted.ID)
	assert.Equal(t, "testuser", formatted.Username)
	assert.Equal(t, "test@example.com", formatted.Email)
	assert.Equal(t, "Test User", formatted.FullName)
	assert.Equal(t, "role123", formatted.RoleID)
	assert.True(t, formatted.IsActive)
}

// Helper functions implementation for testing
func validateUserInput(user model.User) bool {
	if utils.IsEmptyString(user.Username) {
		return false
	}
	if !utils.ValidateEmail(user.Email) {
		return false
	}
	if utils.IsEmptyString(user.FullName) {
		return false
	}
	if utils.IsEmptyString(user.RoleID) {
		return false
	}
	return true
}

func calculatePagination(page, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func calculateTotalPages(total, limit int) int {
	if total == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

func formatUserResponse(user *model.User) *model.User {
	if user == nil {
		return nil
	}
	
	// Create copy without password hash
	formatted := *user
	formatted.PasswordHash = ""
	return &formatted
}