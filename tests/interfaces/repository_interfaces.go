package interfaces

import (
	"uas_backend/app/model"

	"go.mongodb.org/mongo-driver/bson"
)

// UserRepositoryInterface defines the interface for user repository
type UserRepositoryInterface interface {
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id string) (bool, error)
	FindByID(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll() ([]model.User, error)
	GetUserPermissions(userID string) ([]string, error)
}

// AchievementRepositoryInterface defines the interface for achievement repository
type AchievementRepositoryInterface interface {
	Create(achievement *model.Achievement) (*model.Achievement, error)
	FindByID(id string) (*model.Achievement, error)
	FindByIDActive(id string) (*model.Achievement, error)
	Update(achievement *model.Achievement) (*model.Achievement, error)
	UpdateFields(id string, req *model.UpdateAchievementRequest) error
	SoftDelete(id string, deletedBy string) error
	Submit(id string) error
	Verify(id string, verifiedBy string) error
	Reject(id string, reason string) error
	FindAllWithPagination(filter bson.M, page, limit int, sortBy, sortOrder string) ([]model.Achievement, int, error)
	FindByIDs(ids []string) ([]model.Achievement, error)
	GetHistory(id string) ([]model.AchievementHistory, error)
	AddDocument(id string, document model.AchievementDocument) error
}

// AchievementReferenceRepositoryInterface defines the interface for achievement reference repository
type AchievementReferenceRepositoryInterface interface {
	Create(ref *model.AchievementReference) error
	Update(ref *model.AchievementReference) error
	SoftDelete(mongoID string) error
	GetByStudentIDs(studentIDs []string, limit, offset int) ([]model.AchievementReference, error)
	CountByStudentIDs(studentIDs []string) (int, error)
}

// StudentRepositoryInterface defines the interface for student repository
type StudentRepositoryInterface interface {
	FindByUserID(userID string) (*model.Student, error)
	FindByID(id string) (*model.Student, error)
	FindAllWithPagination(limit, offset int) ([]model.Student, error)
	Count() (int, error)
	UpdateAdvisor(studentID, advisorID string) error
	GetStudentIDsByAdvisorID(advisorID string) ([]string, error)
	GetStudentsByAdvisorID(advisorID string) ([]model.Student, error)
	FindAllLecturers(limit, offset int) ([]model.Lecturer, error)
	CountLecturers() (int, error)
	FindLecturerByID(id string) (*model.Lecturer, error)
}