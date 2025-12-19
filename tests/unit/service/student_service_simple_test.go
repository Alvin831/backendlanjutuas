package service_test

import (
	"testing"
	"uas_backend/app/model"
	"uas_backend/tests/mocks"

	"github.com/stretchr/testify/assert"
)

// Simple tests for student service logic using mocks

func TestStudentRepository_FindByID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	studentID := "student123"
	expectedStudent := &model.Student{
		ID:       studentID,
		UserID:   "user123",
		NIM:      "123456789",
		Name:     "John Doe",
		Email:    "john@example.com",
		Program:  "Computer Science",
		Semester: 5,
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("FindByID", studentID).Return(expectedStudent, nil)

	// Test
	student, err := mockRepo.FindByID(studentID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, studentID, student.ID)
	assert.Equal(t, "123456789", student.NIM)
	assert.Equal(t, "John Doe", student.Name)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_FindByUserID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	userID := "user123"
	expectedStudent := &model.Student{
		ID:       "student123",
		UserID:   userID,
		NIM:      "123456789",
		Name:     "John Doe",
		Email:    "john@example.com",
		Program:  "Computer Science",
		Semester: 5,
		IsActive: true,
	}

	// Setup mock expectation
	mockRepo.On("FindByUserID", userID).Return(expectedStudent, nil)

	// Test
	student, err := mockRepo.FindByUserID(userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, userID, student.UserID)
	assert.Equal(t, "123456789", student.NIM)
	assert.Equal(t, "John Doe", student.Name)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_FindAllWithPagination(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	expectedStudents := []model.Student{
		{
			ID:       "student1",
			UserID:   "user1",
			NIM:      "123456789",
			Name:     "John Doe",
			Email:    "john@example.com",
			Program:  "Computer Science",
			Semester: 5,
			IsActive: true,
		},
		{
			ID:       "student2",
			UserID:   "user2",
			NIM:      "987654321",
			Name:     "Jane Smith",
			Email:    "jane@example.com",
			Program:  "Information Technology",
			Semester: 3,
			IsActive: true,
		},
	}

	// Setup mock expectation
	mockRepo.On("FindAllWithPagination", 10, 0).Return(expectedStudents, nil)

	// Test
	students, err := mockRepo.FindAllWithPagination(10, 0)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, students, 2)
	assert.Equal(t, "123456789", students[0].NIM)
	assert.Equal(t, "987654321", students[1].NIM)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_Count(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	expectedCount := 25

	// Setup mock expectation
	mockRepo.On("Count").Return(expectedCount, nil)

	// Test
	count, err := mockRepo.Count()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_UpdateAdvisor(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	studentID := "student123"
	advisorID := "advisor123"

	// Setup mock expectation
	mockRepo.On("UpdateAdvisor", studentID, advisorID).Return(nil)

	// Test
	err := mockRepo.UpdateAdvisor(studentID, advisorID)

	// Assertions
	assert.NoError(t, err)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_GetStudentsByAdvisorID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	advisorID := "advisor123"
	expectedStudents := []model.Student{
		{
			ID:        "student1",
			UserID:    "user1",
			NIM:       "123456789",
			Name:      "John Doe",
			Email:     "john@example.com",
			Program:   "Computer Science",
			Semester:  5,
			AdvisorID: &advisorID,
			IsActive:  true,
		},
		{
			ID:        "student2",
			UserID:    "user2",
			NIM:       "987654321",
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Program:   "Computer Science",
			Semester:  3,
			AdvisorID: &advisorID,
			IsActive:  true,
		},
	}

	// Setup mock expectation
	mockRepo.On("GetStudentsByAdvisorID", advisorID).Return(expectedStudents, nil)

	// Test
	students, err := mockRepo.GetStudentsByAdvisorID(advisorID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, students, 2)
	assert.Equal(t, advisorID, *students[0].AdvisorID)
	assert.Equal(t, advisorID, *students[1].AdvisorID)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_GetStudentIDsByAdvisorID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	advisorID := "advisor123"
	expectedStudentIDs := []string{"student1", "student2", "student3"}

	// Setup mock expectation
	mockRepo.On("GetStudentIDsByAdvisorID", advisorID).Return(expectedStudentIDs, nil)

	// Test
	studentIDs, err := mockRepo.GetStudentIDsByAdvisorID(advisorID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, studentIDs, 3)
	assert.Contains(t, studentIDs, "student1")
	assert.Contains(t, studentIDs, "student2")
	assert.Contains(t, studentIDs, "student3")

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_FindAllLecturers(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	expectedLecturers := []model.Lecturer{
		{
			ID:         "lecturer1",
			UserID:     "user1",
			LecturerID: "NIDN001",
			Name:       "Dr. John Smith",
			Email:      "john.smith@university.edu",
			Department: "Computer Science",
			IsActive:   true,
		},
		{
			ID:         "lecturer2",
			UserID:     "user2",
			LecturerID: "NIDN002",
			Name:       "Dr. Jane Doe",
			Email:      "jane.doe@university.edu",
			Department: "Information Technology",
			IsActive:   true,
		},
	}

	// Setup mock expectation
	mockRepo.On("FindAllLecturers", 10, 0).Return(expectedLecturers, nil)

	// Test
	lecturers, err := mockRepo.FindAllLecturers(10, 0)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, lecturers, 2)
	assert.Equal(t, "NIDN001", lecturers[0].LecturerID)
	assert.Equal(t, "NIDN002", lecturers[1].LecturerID)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_FindLecturerByID(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	lecturerID := "lecturer123"
	expectedLecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     "user123",
		LecturerID: "NIDN001",
		Name:       "Dr. John Smith",
		Email:      "john.smith@university.edu",
		Department: "Computer Science",
		IsActive:   true,
	}

	// Setup mock expectation
	mockRepo.On("FindLecturerByID", lecturerID).Return(expectedLecturer, nil)

	// Test
	lecturer, err := mockRepo.FindLecturerByID(lecturerID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, lecturer)
	assert.Equal(t, lecturerID, lecturer.ID)
	assert.Equal(t, "NIDN001", lecturer.LecturerID)
	assert.Equal(t, "Dr. John Smith", lecturer.Name)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStudentRepository_CountLecturers(t *testing.T) {
	// Setup mock
	mockRepo := new(mocks.MockStudentRepository)

	// Test data
	expectedCount := 15

	// Setup mock expectation
	mockRepo.On("CountLecturers").Return(expectedCount, nil)

	// Test
	count, err := mockRepo.CountLecturers()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}