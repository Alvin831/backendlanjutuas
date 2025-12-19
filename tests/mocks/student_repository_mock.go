package mocks

import (
	"uas_backend/app/model"

	"github.com/stretchr/testify/mock"
)

// MockStudentRepository is a mock implementation of StudentRepository
type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) FindByUserID(userID string) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentRepository) FindByID(id string) (*model.Student, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentRepository) FindAllWithPagination(limit, offset int) ([]model.Student, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *MockStudentRepository) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockStudentRepository) UpdateAdvisor(studentID, advisorID string) error {
	args := m.Called(studentID, advisorID)
	return args.Error(0)
}

func (m *MockStudentRepository) GetStudentIDsByAdvisorID(advisorID string) ([]string, error) {
	args := m.Called(advisorID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStudentRepository) GetStudentsByAdvisorID(advisorID string) ([]model.Student, error) {
	args := m.Called(advisorID)
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *MockStudentRepository) FindAllLecturers(limit, offset int) ([]model.Lecturer, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Lecturer), args.Error(1)
}

func (m *MockStudentRepository) CountLecturers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockStudentRepository) FindLecturerByID(id string) (*model.Lecturer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturer), args.Error(1)
}