package alcs

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Organization struct {
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Notes     sql.NullString `db:"notes"`
	IsActive  bool           `db:"is_active"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (o *Organization) String() string {
	return fmt.Sprintf("Organization{id: %v, name: %v}", o.ID, o.Name)
}

// Organization actions
type Organizations interface {

	// GetAll fetches all organizations
	GetAll() ([]Organization, error)

	// Create instantiates new organization and saves it
	Create(org *Organization) error

	// GetById fetches organization with specific id
	GetById(id int) (*Organization, error)
}

type organizations struct {
	db *sqlx.DB
}

func NewOrganizations(db *sqlx.DB) Organizations {
	return &organizations{db: db}
}

func (o *organizations) GetAll() ([]Organization, error) {
	var orns []Organization
	err := o.db.Select(&orns, "SELECT * FROM organizations")
	return orns, err
}

func (o *organizations) Create(org *Organization) error {
	query := `INSERT INTO organizations (name, notes, is_active, created_at, updated_at)
	          VALUES ($1, $2, true, NOW(), NOW()) RETURNING id`
	return o.db.QueryRow(query, org.Name, org.Notes).Scan(&org.ID)
}

func (o *organizations) GetById(id int) (*Organization, error) {
	var organization Organization
	err := o.db.Get(&organization, "SELECT * FROM organizations WHERE id = $1 LIMIT 1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	}
	return &organization, err
}
