package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/tikhonp/alcs/internal/db/utils"
)

type User struct {
	Id          int            `db:"id"`
	Password    sql.NullString `db:"password"`
	LastLogin   sql.NullTime   `db:"last_login"`
	IsSuperuser bool           `db:"is_superuser"`
	Email       string         `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	IsActive    bool           `db:"is_active"`
	DateJoined  time.Time      `db:"date_joined"`
	OAuthId     sql.NullString `db:"oauth_id"`
}

func (u *User) String() string {
	firstName := "nil"
	if u.FirstName.Valid {
		firstName = u.FirstName.String
	}
	lastName := "nil"
	if u.LastName.Valid {
		lastName = u.LastName.String
	}
	return fmt.Sprintf("User{id: %v, email: %v, firstName: %v, lastName: %v}", u.Id, u.Email, firstName, lastName)
}

// User Actions
type Users interface {

	// GetById returns a user by id
	GetById(id int) (*User, error)

	// IsUserHasPermissions checks if user has specified permissions
	IsUserHasPermissions(userId int, permissionCodenames ...string) (bool, error)

	// FromOAuth finds user by goth user or creates new one
	FromOAuth(guser *goth.User) (*User, error)
}

type users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) Users {
	return &users{db: db}
}

func (u *users) GetById(id int) (*User, error) {
	var user User
	err := u.db.Get(&user, `SELECT * FROM auth_user WHERE id = $1`, id)
	return &user, err
}

func (u *users) IsUserHasPermissions(userId int, permissionCodenames ...string) (bool, error) {
	panic("Not implemented")
}

func (u *users) FromOAuth(guser *goth.User) (*User, error) {
	var user User
	err := u.db.Get(
		&user,
		`SELECT * FROM auth_user WHERE oauth_id = $1 LIMIT 1`,
		guser.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		// Create new user
		err = u.db.Get(
			&user,
			`INSERT INTO auth_user (oauth_id, email, first_name, last_name, date_joined) VALUES ($1, $2, $3, $4, $5) RETURNING *`,
			guser.UserID,
			guser.Email,
			utils.NewNullString(guser.FirstName),
            utils.NewNullString(guser.LastName),
			time.Now(),
		)
		return &user, err
	}

	return &user, err
}
