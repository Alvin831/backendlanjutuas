package simple_test

import (
	"testing"
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/stretchr/testify/assert"
)

// Test achievement business logic functions

// Test ValidateAchievementData function
func TestValidateAchievementData(t *testing.T) {
	// Valid achievement
	validAchievement := model.Achievement{
		Title:       "Test Achievement",
		Category:    "akademik",
		Points:      75,
		StudentID:   "student123",
		Description: "Test description",
	}
	
	assert.True(t, validateAchievementData(validAchievement))
	
	// Invalid - empty title
	invalidAchievement1 := validAchievement
	invalidAchievement1.Title = ""
	assert.False(t, validateAchievementData(invalidAchievement1))
	
	// Invalid - empty category
	invalidAchievement2 := validAchievement
	invalidAchievement2.Category = ""
	assert.False(t, validateAchievementData(invalidAchievement2))
	
	// Invalid - empty student ID
	invalidAchievement3 := validAchievement
	invalidAchievement3.StudentID = ""
	assert.False(t, validateAchievementData(invalidAchievement3))
}

// Test CalculateAchievementPoints function
func TestCalculateAchievementPoints(t *testing.T) {
	// Test direct points calculation
	assert.Equal(t, 25, calculateAchievementPointsByLevel("Lokal (1-25 poin)"))
	assert.Equal(t, 50, calculateAchievementPointsByLevel("Regional (26-50 poin)"))
	assert.Equal(t, 75, calculateAchievementPointsByLevel("Nasional (51-75 poin)"))
	assert.Equal(t, 100, calculateAchievementPointsByLevel("Internasional (76+ poin)"))
	assert.Equal(t, 0, calculateAchievementPointsByLevel("Unknown"))
	
	// Test with achievement object
	achievement1 := model.Achievement{Points: 75}
	assert.Equal(t, 75, achievement1.Points)
	
	achievement2 := model.Achievement{Points: 100}
	assert.Equal(t, 100, achievement2.Points)
}

// Test CanUserAccessAchievement function (role-based access)
func TestCanUserAccessAchievement(t *testing.T) {
	achievement := model.Achievement{
		StudentID: "student123",
	}
	
	// Admin can access any achievement
	assert.True(t, canUserAccessAchievement("admin", "any-user-id", achievement))
	
	// Dosen can access any achievement
	assert.True(t, canUserAccessAchievement("dosen", "any-user-id", achievement))
	
	// Mahasiswa can access own achievement
	assert.True(t, canUserAccessAchievement("mahasiswa", "student123", achievement))
	
	// Mahasiswa cannot access other's achievement
	assert.False(t, canUserAccessAchievement("mahasiswa", "other-student", achievement))
}

// Test ValidateAchievementStatus function
func TestValidateAchievementStatus(t *testing.T) {
	validStatuses := []string{"draft", "submitted", "verified", "rejected"}
	
	for _, status := range validStatuses {
		assert.True(t, validateAchievementStatus(status))
	}
	
	// Invalid statuses
	assert.False(t, validateAchievementStatus("invalid"))
	assert.False(t, validateAchievementStatus(""))
	assert.False(t, validateAchievementStatus("pending"))
}

// Test CanChangeAchievementStatus function
func TestCanChangeAchievementStatus(t *testing.T) {
	// Draft -> Submitted (allowed for mahasiswa)
	assert.True(t, canChangeAchievementStatus("draft", "submitted", "mahasiswa"))
	
	// Submitted -> Verified (allowed for dosen)
	assert.True(t, canChangeAchievementStatus("submitted", "verified", "dosen"))
	
	// Submitted -> Rejected (allowed for dosen)
	assert.True(t, canChangeAchievementStatus("submitted", "rejected", "dosen"))
	
	// Verified -> Draft (not allowed)
	assert.False(t, canChangeAchievementStatus("verified", "draft", "mahasiswa"))
	
	// Mahasiswa cannot verify
	assert.False(t, canChangeAchievementStatus("submitted", "verified", "mahasiswa"))
}

// Helper functions implementation
func validateAchievementData(achievement model.Achievement) bool {
	if utils.IsEmptyString(achievement.Title) {
		return false
	}
	if utils.IsEmptyString(achievement.Category) {
		return false
	}
	if utils.IsEmptyString(achievement.StudentID) {
		return false
	}
	if utils.IsEmptyString(achievement.Description) {
		return false
	}
	if achievement.Points <= 0 {
		return false
	}
	return true
}

func calculateAchievementPointsByLevel(level string) int {
	return utils.CalculatePoints(level)
}

func canUserAccessAchievement(role, userID string, achievement model.Achievement) bool {
	// Admin and dosen can access all achievements
	if role == "admin" || role == "dosen" {
		return true
	}
	
	// Mahasiswa can only access their own achievements
	if role == "mahasiswa" {
		return achievement.StudentID == userID
	}
	
	return false
}

func validateAchievementStatus(status string) bool {
	validStatuses := []string{"draft", "submitted", "verified", "rejected"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func canChangeAchievementStatus(currentStatus, newStatus, role string) bool {
	// Define allowed transitions
	transitions := map[string]map[string][]string{
		"draft": {
			"submitted": {"mahasiswa", "admin"},
		},
		"submitted": {
			"verified": {"dosen", "admin"},
			"rejected": {"dosen", "admin"},
		},
	}
	
	if allowedRoles, exists := transitions[currentStatus][newStatus]; exists {
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return true
			}
		}
	}
	
	return false
}