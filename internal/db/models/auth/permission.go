package auth

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Permission struct {
	name string
	code string
}

// All Rights and adding new organizations
var SuperAdmin = Permission{name: "All Rights and adding new organizations", code: "superadmin"}

type Permissions interface {

	// AddPermissionForUser creates specified permission if it is not exists
	// and creates relation between user and permission.
	//
	// WARN: userId must exist.
	AddPermissionForUser(userId int, perm Permission) error

	// RemovePermissionForUser removes relation between permission and user.
	RemovePermissionForUser(userId int, perm Permission) error
}

type permissions struct {
	db *sqlx.DB
}

func NewPermissions(db *sqlx.DB) Permissions {
	return &permissions{db: db}
}

func (p *permissions) AddPermissionForUser(userId int, perm Permission) error {
	var permissionId int
	// get permission id, check if it exists
	rows, err := p.db.Query("SELECT id FROM auth_permission WHERE codename = $1", perm.code)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&permissionId)
	} else {
        rows, err = p.db.Query("INSERT INTO auth_permission (name, codename) VALUES ($1, $2) RETURNING id", perm.name, perm.code)
        if err != nil {
            return err
        }
        defer rows.Close()
        if rows.Next() {
            rows.Scan(&permissionId)
        } else {
            return errors.New("no id in insert permission statements")
        }
	}
	_, err = p.db.Query("INSERT INTO auth_user_permissions (user_id, permission_id) VALUES ($1, $2)", userId, permissionId)
	return err
}

func (p *permissions) RemovePermissionForUser(userId int, perm Permission) error {
	panic("Not implemented")
}
