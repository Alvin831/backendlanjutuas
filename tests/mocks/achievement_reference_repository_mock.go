package mocks

import (
	"uas_backend/app/model"

	"github.com/stretchr/testify/mock"
)

// MockAchievementReferenceRepository is a mock implementation of AchievementReferenceRepository
type MockAchievementReferenceRepository struct {
	mock.Mock
}

func (m *MockAchievementReferenceRepository) Create(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockAchievementReferenceRepository) Update(ref *model.AchievementReference) error {
	args := m.Called(ref)
	return args.Error(0)
}

func (m *MockAchievementReferenceRepository) SoftDelete(mongoID string) error {
	args := m.Called(mongoID)
	return args.Error(0)
}

func (m *MockAchievementReferenceRepository) GetByStudentIDs(studentIDs []string, limit, offset int) ([]model.AchievementReference, error) {
	args := m.Called(studentIDs, limit, offset)
	return args.Get(0).([]model.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) CountByStudentIDs(studentIDs []string) (int, error) {
	args := m.Called(studentIDs)
	return args.Int(0), args.Error(1)
}