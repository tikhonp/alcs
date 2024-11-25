package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/tikhonp/alcs/internal/db/models/auth"
)

type ModelsFactory interface {
	AuthGroups() auth.Groups
	AuthPermissions() auth.Permissions
	AuthUsers() auth.Users
}

type modelsFactory struct {
	users auth.Users
}

func newModelsFactory(db *sqlx.DB) ModelsFactory {
	return &modelsFactory{
		users: auth.NewUsers(db),
	}
}

func (f *modelsFactory) AuthGroups() auth.Groups {
	panic("not implemented")
}

func (f *modelsFactory) AuthPermissions() auth.Permissions {
	panic("not implemented")
}

func (f *modelsFactory) AuthUsers() auth.Users {
	return f.users
}
