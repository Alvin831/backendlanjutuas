package repository

import (
	"database/sql"
	"uas_backend/app/model"
)

type PermissionRepository struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetAll() ([]model.Permission, error) {
	rows, err := r.db.Query(`SELECT id, name, resource, action, description FROM permissions ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var p model.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindByID(id string) (*model.Permission, error) {
	var p model.Permission
	err := r.db.QueryRow(
		`SELECT id, name, resource, action, description FROM permissions WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *PermissionRepository) Create(p *model.Permission) (*model.Permission, error) {
	err := r.db.QueryRow(`
		INSERT INTO permissions (name, resource, action, description)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		p.Name, p.Resource, p.Action, p.Description,
	).Scan(&p.ID)

	return p, err
}

func (r *PermissionRepository) Delete(id string) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM permissions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// =================== ROLE – PERMISSIONS =====================

func (r *PermissionRepository) Assign(roleID string, permissionID string) error {
	_, err := r.db.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)`, roleID, permissionID)
	return err
}

func (r *PermissionRepository) Remove(roleID string, permissionID string) error {
	_, err := r.db.Exec(`
		DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`,
		roleID, permissionID)
	return err
}
