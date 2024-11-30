package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/tikhonp/alcs/internal/db/models/alcs"
	"github.com/tikhonp/alcs/internal/db/models/auth"
)

type ModelsFactory interface {
	AuthUsers() auth.Users
	AlcsOrganizations() alcs.Organizations
}

type modelsFactory struct {
	users         auth.Users
	organizations alcs.Organizations
}

func newModelsFactory(db *sqlx.DB) ModelsFactory {
	return &modelsFactory{
		users:         auth.NewUsers(db),
		organizations: alcs.NewOrganizations(db),
	}
}

func (f *modelsFactory) AuthUsers() auth.Users {
	return f.users
}

func (f *modelsFactory) AlcsOrganizations() alcs.Organizations {
	return f.organizations
}
