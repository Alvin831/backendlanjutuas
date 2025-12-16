package repository

import (
	"database/sql"
	"fmt"
	"time"
	"uas_backend/app/model"
)

type AchievementReferenceRepository struct {
	db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) *AchievementReferenceRepository {
	return &AchievementReferenceRepository{db: db}
}

// Create achievement reference
func (r *AchievementReferenceRepository) Create(ref *model.AchievementReference) error {
	query := `
		INSERT INTO achievement_references (mongo_achievement_id, student_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	err := r.db.QueryRow(query,
		ref.MongoID, ref.StudentID, ref.Status, ref.CreatedAt, ref.UpdatedAt,
	).Scan(&ref.ID)
	
	return err
}

// Update achievement reference
func (r *AchievementReferenceRepository) Update(ref *model.AchievementReference) error {
	query := `
		UPDATE achievement_references 
		SET status=$1, updated_at=$2
		WHERE mongo_achievement_id=$3
	`
	
	_, err := r.db.Exec(query, ref.Status, ref.UpdatedAt, ref.MongoID)
	return err
}

// Soft delete achievement reference
func (r *AchievementReferenceRepository) SoftDelete(mongoID string) error {
	query := `
		UPDATE achievement_references 
		SET is_deleted=true, deleted_at=$1, updated_at=$2
		WHERE mongo_achievement_id=$3
	`
	
	now := time.Now()
	_, err := r.db.Exec(query, now, now, mongoID)
	return err
}

// Find by mongo ID
func (r *AchievementReferenceRepository) FindByMongoID(mongoID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, mongo_achievement_id, student_id, status, created_at, updated_at
		FROM achievement_references 
		WHERE mongo_achievement_id=$1
	`
	
	ref := &model.AchievementReference{}
	err := r.db.QueryRow(query, mongoID).Scan(
		&ref.ID, &ref.MongoID, &ref.StudentID, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return ref, err
}

// Get achievement references by student IDs with pagination (FR-006)
func (r *AchievementReferenceRepository) GetByStudentIDs(studentIDs []string, limit, offset int) ([]model.AchievementReference, error) {
	if len(studentIDs) == 0 {
		return []model.AchievementReference{}, nil
	}

	// Convert user_ids to student_ids by joining with students table
	placeholders := ""
	args := make([]interface{}, 0, len(studentIDs)+2)
	
	for i, studentID := range studentIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + fmt.Sprintf("%d", i+1)
		args = append(args, studentID)
	}
	
	query := fmt.Sprintf(`
		SELECT ar.id, ar.mongo_achievement_id, ar.student_id, ar.status, ar.created_at, ar.updated_at
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE s.user_id IN (%s) AND (ar.is_deleted = false OR ar.is_deleted IS NULL)
		ORDER BY ar.created_at DESC
		LIMIT $%d OFFSET $%d
	`, placeholders, len(studentIDs)+1, len(studentIDs)+2)
	
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID, &ref.MongoID, &ref.StudentID, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

// Count achievement references by student IDs
func (r *AchievementReferenceRepository) CountByStudentIDs(studentIDs []string) (int, error) {
	if len(studentIDs) == 0 {
		return 0, nil
	}

	placeholders := ""
	args := make([]interface{}, 0, len(studentIDs))
	
	for i, studentID := range studentIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + fmt.Sprintf("%d", i+1)
		args = append(args, studentID)
	}
	
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM achievement_references ar
		JOIN students s ON ar.student_id = s.id
		WHERE s.user_id IN (%s) AND (ar.is_deleted = false OR ar.is_deleted IS NULL)
	`, placeholders)
	
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}