package repository

import (
	"database/sql"
	"uas_backend/app/model"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive, time.Now(), time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	_, err := r.db.Exec(`
		UPDATE users
		SET username=$1, email=$2, full_name=$3, role_id=$4, is_active=$5, updated_at=$6
		WHERE id=$7
	`, user.Username, user.Email, user.FullName, user.RoleID, user.IsActive, time.Now(), user.ID)
	if err != nil {
		return nil, err
	}
	return r.FindByID(user.ID)
}

func (r *UserRepository) Delete(id string) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	ra, _ := res.RowsAffected()
	return ra > 0, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}
	row := r.db.QueryRow(`
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, id)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	u := &model.User{}
	row := r.db.QueryRow(`
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`, username)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindAll() ([]model.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

// GetUserPermissions - Ambil permissions berdasarkan role user
func (r *UserRepository) GetUserPermissions(userID string) ([]string, error) {
	query := `
		SELECT DISTINCT p.name
		FROM users u
		JOIN roles ro ON u.role_id = ro.id
		JOIN role_permissions rp ON ro.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE u.id = $1 AND u.is_active = true
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	
	return permissions, nil
}
