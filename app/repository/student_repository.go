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

// Find all students with pagination
func (r *StudentRepository) FindAllWithPagination(limit, offset int) ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at,
		       u.full_name, u.email, u.is_active
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.student_id ASC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var student model.Student
		var fullName, email sql.NullString
		var isActive sql.NullBool
		
		err := rows.Scan(
			&student.ID, &student.UserID, &student.NIM, &student.Program, 
			&student.Semester, &student.AdvisorID, &student.CreatedAt,
			&fullName, &email, &isActive,
		)
		if err != nil {
			return nil, err
		}
		
		// Set values from user table
		if fullName.Valid {
			student.Name = fullName.String
		} else {
			student.Name = student.NIM
		}
		if email.Valid {
			student.Email = email.String
		}
		if isActive.Valid {
			student.IsActive = isActive.Bool
		} else {
			student.IsActive = true
		}
		
		students = append(students, student)
	}

	return students, nil
}

// Count total students
func (r *StudentRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM students").Scan(&count)
	return count, err
}

// Find student by ID
func (r *StudentRepository) FindByID(id string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at,
		       u.full_name, u.email, u.is_active
		FROM students s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = $1
	`
	
	student := &model.Student{}
	var fullName, email sql.NullString
	var isActive sql.NullBool
	
	err := r.db.QueryRow(query, id).Scan(
		&student.ID, &student.UserID, &student.NIM, &student.Program, 
		&student.Semester, &student.AdvisorID, &student.CreatedAt,
		&fullName, &email, &isActive,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Set values from user table
	if fullName.Valid {
		student.Name = fullName.String
	} else {
		student.Name = student.NIM
	}
	if email.Valid {
		student.Email = email.String
	}
	if isActive.Valid {
		student.IsActive = isActive.Bool
	} else {
		student.IsActive = true
	}
	
	return student, nil
}

// Update student advisor
func (r *StudentRepository) UpdateAdvisor(studentID, advisorID string) error {
	query := `UPDATE students SET advisor_id = $1 WHERE id = $2`
	_, err := r.db.Exec(query, advisorID, studentID)
	return err
}

// Find all lecturers with pagination
func (r *StudentRepository) FindAllLecturers(limit, offset int) ([]model.Lecturer, error) {
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at,
		       u.full_name, u.email, u.is_active
		FROM lecturers l
		LEFT JOIN users u ON l.user_id = u.id
		ORDER BY l.lecturer_id ASC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []model.Lecturer
	for rows.Next() {
		var lecturer model.Lecturer
		var fullName, email sql.NullString
		var isActive sql.NullBool
		
		err := rows.Scan(
			&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt,
			&fullName, &email, &isActive,
		)
		if err != nil {
			return nil, err
		}
		
		// Set values from user table
		if fullName.Valid {
			lecturer.Name = fullName.String
		} else {
			lecturer.Name = lecturer.LecturerID
		}
		if email.Valid {
			lecturer.Email = email.String
		}
		if isActive.Valid {
			lecturer.IsActive = isActive.Bool
		} else {
			lecturer.IsActive = true
		}
		
		lecturers = append(lecturers, lecturer)
	}

	return lecturers, nil
}

// Count total lecturers
func (r *StudentRepository) CountLecturers() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM lecturers").Scan(&count)
	return count, err
}

// Find lecturer by ID
func (r *StudentRepository) FindLecturerByID(id string) (*model.Lecturer, error) {
	query := `
		SELECT l.id, l.user_id, l.lecturer_id, l.department, l.created_at,
		       u.full_name, u.email, u.is_active
		FROM lecturers l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE l.id = $1
	`
	
	lecturer := &model.Lecturer{}
	var fullName, email sql.NullString
	var isActive sql.NullBool
	
	err := r.db.QueryRow(query, id).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt,
		&fullName, &email, &isActive,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Set values from user table
	if fullName.Valid {
		lecturer.Name = fullName.String
	} else {
		lecturer.Name = lecturer.LecturerID
	}
	if email.Valid {
		lecturer.Email = email.String
	}
	if isActive.Valid {
		lecturer.IsActive = isActive.Bool
	} else {
		lecturer.IsActive = true
	}
	
	return lecturer, nil
}