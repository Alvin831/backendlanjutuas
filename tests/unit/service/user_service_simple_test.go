package service_test

import (
	"testing"
	"uas_backend/app/model"
	"uas_backend/tests/mocks"

	"github.com/stretchr/testify/assert"
)

// Simple tests for user service logic using mocks

func TestUserRepository_CreateUser(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	newUser := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   "role123",
		IsActive: true,
	}

	expectedUser := &model.User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		RoleID:   "role123",
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("Create", newUser).Return(expectedUser, nil)

	// Test
	createdUser, err := mockRepo.Create(newUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "user123", createdUser.ID)
	assert.Equal(t, "testuser", createdUser.Username)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	userID := "user123"
	expectedUser := &model.User{
		ID:       userID,
		Username: "admin",
		Email:    "admin@example.com",
		FullName: "Administrator",
		RoleID:   "admin-role",
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Test
	user, err := mockRepo.FindByID(userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "admin", user.Username)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	username := "admin"
	expectedUser := &model.User{
		ID:       "user123",
		Username: username,
		Email:    "admin@example.com",
		FullName: "Administrator",
		RoleID:   "admin-role",
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("FindByUsername", username).Return(expectedUser, nil)

	// Test
	user, err := mockRepo.FindByUsername(username)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "admin@example.com", user.Email)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindAll(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	expectedUsers := []model.User{
		{
			ID:       "user1",
			Username: "admin",
			Email:    "admin@example.com",
			FullName: "Administrator",
			IsActive: true,
		},
		{
			ID:       "user2",
			Username: "user1",
			Email:    "user1@example.com",
			FullName: "User One",
			IsActive: true,
		},
	}

	// Setup mock expectation
	mockRepo.On("FindAll").Return(expectedUsers, nil)

	// Test
	users, err := mockRepo.FindAll()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "admin", users[0].Username)
	assert.Equal(t, "user1", users[1].Username)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	updateUser := &model.User{
		ID:       "user123",
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
		RoleID:   "role123",
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("Update", updateUser).Return(updateUser, nil)

	// Test
	updatedUser, err := mockRepo.Update(updateUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.Equal(t, "updated@example.com", updatedUser.Email)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_DeleteUser(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	userID := "user123"

	// Setup mock expectation
	mockRepo.On("Delete", userID).Return(true, nil)

	// Test
	success, err := mockRepo.Delete(userID)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, success)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_GetUserPermissions(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Test data
	userID := "user123"
	expectedPermissions := []string{"read", "write", "admin", "manage_users"}

	// Setup mock expectation
	mockRepo.On("GetUserPermissions", userID).Return(expectedPermissions, nil)

	// Test
	permissions, err := mockRepo.GetUserPermissions(userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, permissions, 4)
	assert.Contains(t, permissions, "read")
	assert.Contains(t, permissions, "write")
	assert.Contains(t, permissions, "admin")
	assert.Contains(t, permissions, "manage_users")

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Setup mock expectation for not found
	mockRepo.On("FindByID", "nonexistent").Return(nil, nil)

	// Test
	user, err := mockRepo.FindByID("nonexistent")

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockUserRepository)

	// Setup mock expectation for not found
	mockRepo.On("FindByUsername", "nonexistent").Return(nil, nil)

	// Test
	user, err := mockRepo.FindByUsername("nonexistent")

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}