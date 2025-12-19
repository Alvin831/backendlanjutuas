package mocks

import (
	"uas_backend/app/model"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

// MockAchievementRepository is a mock implementation of AchievementRepository
type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) Create(achievement *model.Achievement) (*model.Achievement, error) {
	args := m.Called(achievement)
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) FindByID(id string) (*model.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) FindByIDActive(id string) (*model.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) Update(achievement *model.Achievement) (*model.Achievement, error) {
	args := m.Called(achievement)
	return args.Get(0).(*model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) UpdateFields(id string, req *model.UpdateAchievementRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockAchievementRepository) SoftDelete(id string, deletedBy string) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}

func (m *MockAchievementRepository) Submit(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAchievementRepository) Verify(id string, verifiedBy string) error {
	args := m.Called(id, verifiedBy)
	return args.Error(0)
}

func (m *MockAchievementRepository) Reject(id string, reason string) error {
	args := m.Called(id, reason)
	return args.Error(0)
}

func (m *MockAchievementRepository) FindAllWithPagination(filter bson.M, page, limit int, sortBy, sortOrder string) ([]model.Achievement, int, error) {
	args := m.Called(filter, page, limit, sortBy, sortOrder)
	return args.Get(0).([]model.Achievement), args.Int(1), args.Error(2)
}

func (m *MockAchievementRepository) FindByIDs(ids []string) ([]model.Achievement, error) {
	args := m.Called(ids)
	return args.Get(0).([]model.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) GetHistory(id string) ([]model.AchievementHistory, error) {
	args := m.Called(id)
	return args.Get(0).([]model.AchievementHistory), args.Error(1)
}

func (m *MockAchievementRepository) AddDocument(id string, document model.AchievementDocument) error {
	args := m.Called(id, document)
	return args.Error(0)
}