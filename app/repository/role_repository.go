package repository

import (
	"database/sql"
	"uas_backend/app/model"
	"time"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *model.Role) (*model.Role, error) {
	err := r.db.QueryRow(`
		INSERT INTO roles (name, description, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, role.Name, role.Description, time.Now()).Scan(&role.ID, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetAll() ([]model.Role, error) {
	rows, err := r.db.Query(`SELECT id, name, description, created_at FROM roles ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var ro model.Role
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.Description, &ro.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, ro)
	}
	return roles, nil
}

func (r *RoleRepository) FindByID(id string) (*model.Role, error) {
	var ro model.Role
	err := r.db.QueryRow(`SELECT id, name, description, created_at FROM roles WHERE id = $1`, id).
		Scan(&ro.ID, &ro.Name, &ro.Description, &ro.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *RoleRepository) FindByName(name string) (*model.Role, error) {
	var ro model.Role
	err := r.db.QueryRow(`SELECT id, name, description, created_at FROM roles WHERE name = $1`, name).
		Scan(&ro.ID, &ro.Name, &ro.Description, &ro.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ro, nil
}
