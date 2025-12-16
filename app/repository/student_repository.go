package repository

import (
	"database/sql"
	"uas_backend/app/model"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Get student IDs by advisor ID (FR-006)
func (r *StudentRepository) GetStudentIDsByAdvisorID(advisorID string) ([]string, error) {
	// Get lecturer ID from user_id, then get students
	query := `
		SELECT s.user_id 
		FROM students s
		JOIN lecturers l ON s.advisor_id = l.id
		WHERE l.user_id = $1
	`
	
	rows, err := r.db.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []string
	for rows.Next() {
		var studentID string
		if err := rows.Scan(&studentID); err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, studentID)
	}

	return studentIDs, nil
}

// Get students by advisor ID with details
func (r *StudentRepository) GetStudentsByAdvisorID(advisorID string) ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at
		FROM students s
		JOIN lecturers l ON s.advisor_id = l.id
		WHERE l.user_id = $1
		ORDER BY s.student_id ASC
	`
	
	rows, err := r.db.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(
			&student.ID, &student.UserID, &student.NIM, &student.Program, 
			&student.Semester, &student.AdvisorID, &student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Set default values for missing fields
		student.Name = student.NIM // Use student_id as name for now
		student.Email = ""
		student.IsActive = true
		students = append(students, student)
	}

	return students, nil
}

// Find student by user ID
func (r *StudentRepository) FindByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students 
		WHERE user_id = $1
	`
	
	student := &model.Student{}
	err := r.db.QueryRow(query, userID).Scan(
		&student.ID, &student.UserID, &student.NIM, &student.Program, 
		&student.Semester, &student.AdvisorID, &student.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	// Set default values
	student.Name = student.NIM
	student.Email = ""
	student.IsActive = true
	
	return student, err
}