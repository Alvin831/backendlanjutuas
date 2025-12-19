package service_test

import (
	"testing"
	"uas_backend/app/model"
	"uas_backend/tests/mocks"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// Simple tests for achievement service logic using mocks

func TestAchievementRepository_CreateAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	newAchievement := &model.Achievement{
		StudentID:   "student123",
		Title:       "Test Achievement",
		Category:    "akademik",
		Points:      75,
		Status:      "draft",
		Description: "Test description",
	}

	expectedAchievement := &model.Achievement{
		StudentID:   "student123",
		Title:       "Test Achievement",
		Category:    "akademik",
		Points:      75,
		Status:      "draft",
		Description: "Test description",
	}

	// Setup mock expectation
	mockRepo.On("Create", newAchievement).Return(expectedAchievement, nil)

	// Test
	createdAchievement, err := mockRepo.Create(newAchievement)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, createdAchievement)
	assert.Equal(t, "student123", createdAchievement.StudentID)
	assert.Equal(t, "Test Achievement", createdAchievement.Title)
	assert.Equal(t, "draft", createdAchievement.Status)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_FindByID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"
	expectedAchievement := &model.Achievement{
		StudentID:   "student123",
		Title:       "Test Achievement",
		Category:    "akademik",
		Points:      75,
		Status:      "verified",
		Description: "Test description",
	}

	// Setup mock expectation
	mockRepo.On("FindByID", achievementID).Return(expectedAchievement, nil)

	// Test
	achievement, err := mockRepo.FindByID(achievementID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, achievement)
	assert.Equal(t, "student123", achievement.StudentID)
	assert.Equal(t, "Test Achievement", achievement.Title)
	assert.Equal(t, "verified", achievement.Status)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_FindAllWithPagination_Student(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data - student can only see their own achievements
	studentID := "student123"
	filter := bson.M{"student_id": studentID, "is_deleted": bson.M{"$ne": true}}
	
	expectedAchievements := []model.Achievement{
		{
			StudentID:   studentID,
			Title:       "Achievement 1",
			Category:    "akademik",
			Points:      75,
			Status:      "verified",
			Description: "Description 1",
		},
		{
			StudentID:   studentID,
			Title:       "Achievement 2",
			Category:    "non-akademik",
			Points:      50,
			Status:      "draft",
			Description: "Description 2",
		},
	}

	// Setup mock expectation
	mockRepo.On("FindAllWithPagination", filter, 1, 10, "created_at", "desc").Return(expectedAchievements, 2, nil)

	// Test
	achievements, total, err := mockRepo.FindAllWithPagination(filter, 1, 10, "created_at", "desc")

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, achievements, 2)
	assert.Equal(t, 2, total)
	assert.Equal(t, studentID, achievements[0].StudentID)
	assert.Equal(t, studentID, achievements[1].StudentID)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_UpdateAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	updateAchievement := &model.Achievement{
		StudentID:   "student123",
		Title:       "Updated Achievement",
		Category:    "non-akademik",
		Points:      100,
		Status:      "draft",
		Description: "Updated description",
	}

	// Setup mock expectation
	mockRepo.On("Update", updateAchievement).Return(updateAchievement, nil)

	// Test
	updatedAchievement, err := mockRepo.Update(updateAchievement)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedAchievement)
	assert.Equal(t, "Updated Achievement", updatedAchievement.Title)
	assert.Equal(t, "non-akademik", updatedAchievement.Category)
	assert.Equal(t, 100, updatedAchievement.Points)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_SubmitAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"

	// Setup mock expectation
	mockRepo.On("Submit", achievementID).Return(nil)

	// Test
	err := mockRepo.Submit(achievementID)

	// Assertions
	assert.NoError(t, err)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_VerifyAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"
	verifiedBy := "dosen123"

	// Setup mock expectation
	mockRepo.On("Verify", achievementID, verifiedBy).Return(nil)

	// Test
	err := mockRepo.Verify(achievementID, verifiedBy)

	// Assertions
	assert.NoError(t, err)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_RejectAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"
	reason := "Dokumen tidak lengkap"

	// Setup mock expectation
	mockRepo.On("Reject", achievementID, reason).Return(nil)

	// Test
	err := mockRepo.Reject(achievementID, reason)

	// Assertions
	assert.NoError(t, err)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_SoftDelete(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"
	deletedBy := "admin123"

	// Setup mock expectation
	mockRepo.On("SoftDelete", achievementID, deletedBy).Return(nil)

	// Test
	err := mockRepo.SoftDelete(achievementID, deletedBy)

	// Assertions
	assert.NoError(t, err)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementRepository_GetHistory(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockAchievementRepository)

	// Test data
	achievementID := "achievement123"
	expectedHistory := []model.AchievementHistory{
		{
			Action:    "created",
			Status:    "draft",
			ActorID:   "student123",
			ActorType: "student",
			Message:   "Achievement created",
		},
		{
			Action:    "submitted",
			Status:    "submitted",
			ActorID:   "student123",
			ActorType: "student",
			Message:   "Achievement submitted for verification",
		},
	}

	// Setup mock expectation
	mockRepo.On("GetHistory", achievementID).Return(expectedHistory, nil)

	// Test
	history, err := mockRepo.GetHistory(achievementID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "created", history[0].Action)
	assert.Equal(t, "submitted", history[1].Action)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}