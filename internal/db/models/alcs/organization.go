package alcs

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Organization struct {
	Id    int            `db:"id"`
	Name  string         `db:"name"`
	Notes sql.NullString `db:"notes"`
}

func (o *Organization) String() string {
	return fmt.Sprintf("Organization{id: %v, name: %v}", o.Id, o.Name)
}

// Organization actions
type Organizations interface {

	// GetAll fetches all organizations
	GetAll() ([]Organization, error)

	// Create instantiates new organization and saves it
	Create(name, notes string) error
}

type organizations struct {
	db *sqlx.DB
}

func NewOrganizations(db *sqlx.DB) Organizations {
	return &organizations{db: db}
}

func (o *organizations) GetAll() ([]Organization, error) {
	var orns []Organization
	err := o.db.Select(&orns, "SELECT * FROM alcs_organization")
	return orns, err
}

func (o *organizations) Create(name, notes string) error {
	panic("not implemented")
}
